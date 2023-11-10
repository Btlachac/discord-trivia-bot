package worker

import (
	"context"
	"encoding/binary"
	"fmt"
	"go-trivia-api/internal/db"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	lightningRoundType = "LIGHTNING"
	listRoundType      = "LIST"

	listRoundInstructions = "This is a List Round! It's one question with multiple correct answers... teams get 1 point for each correct answer they write down."
	megaRoundReminder     = `**Reminder about the Mega Round:**\n
	You can pick any regular trivia round as your mega round, **excluding** the **list round**, the **image round** and the **audio round**.
	 To do so you number your answers from 5-1 and you will get that many points if that answer is correct.`
)

type triviaService interface {
	GetNewTrivia(ctx context.Context) (db.Trivia, error)
	MarkTriviaUsed(ctx context.Context, triviaId int64) error
}

// TODO: add mutex lock around trivia in progress
// TODO: do validation on channel IDs in main
type Bot struct {
	triviaHost    *discordgo.Session
	triviaService triviaService

	triviaChannelID  string
	commandChannelID string

	imageRoundSleepDelay time.Duration
	roundSleepDelay      time.Duration
	questionSleepDelay   time.Duration
	roundStartSleepDelay time.Duration

	currentTrivia db.Trivia
	triviaLock    sync.Mutex
	audioBuffer   [][]byte
}

func NewBot(
	triviaHost *discordgo.Session,
	triviaService triviaService,
	triviaChannelID string,
	commandChannelID string,
	imageRoundSleepDelay time.Duration,
	roundSleepDelay time.Duration,
	questionSleepDelay time.Duration,
	roundStartSleepDelay time.Duration,
) *Bot {
	return &Bot{
		triviaHost:           triviaHost,
		triviaService:        triviaService,
		triviaChannelID:      triviaChannelID,
		commandChannelID:     commandChannelID,
		imageRoundSleepDelay: imageRoundSleepDelay,
		roundSleepDelay:      roundSleepDelay,
		questionSleepDelay:   questionSleepDelay,
		roundStartSleepDelay: roundStartSleepDelay,
		audioBuffer:          make([][]byte, 0),
	}
}

func (b *Bot) Run() error {
	b.triviaHost.AddHandler(b.messageCreateHandler)

	//TODO: not sure we need this?
	b.triviaHost.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	err := b.triviaHost.Open()
	if err != nil {
		return fmt.Errorf("Error opening Discord session: %w", err)
	}
	defer b.triviaHost.Close()

	slog.Info("Trivia Bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return nil
}

// general func to listen for messages
func (b *Bot) messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore messages not in the bot commands channel
	if m.ChannelID != b.commandChannelID {
		return
	}

	message := m.Content

	switch message {
	case "!help":
		//TODO
	case "!start":
		if err := b.startTrivia(); err != nil {
			slog.Error("error running trivia", err)
			//TODO: send message to channel about error
		}
	case "!stop":
		//TODO
	case "!answer",
		"!answers":
		//TODO
	case "!next":
		if err := b.getNextTrivia(); err != nil {
			slog.Error("eror fetching trivia", err)
		}
	case "!audio":
		if err := b.playAudio(); err != nil {
			slog.Error("error running audio round", err)
		}
	case "!image":
		//TODO
	default:
	}
}

// TODO: think about if we want to return for any error? If we miss 1 message it might be better to just log it? Perhaps create a budget,
// if more than 5 messages fail we die
func (b *Bot) startTrivia() error {
	if !b.triviaLock.TryLock() {
		// TODO send message to botcommands to signify trivia in progress
		return fmt.Errorf("trivia already in progress")
	}
	defer b.triviaLock.Unlock()
	var (
		err    error
		trivia db.Trivia
	)

	//Fetch Trivia
	err = b.getNextTrivia()
	if err != nil {
		return err
	}

	trivia = b.currentTrivia

	// Send Trivia overview
	err = b.hostMessageTriviaChannel("**Trivia Overview:**")
	if err != nil {
		return err
	}

	//TODO: move to func
	overviewMessage := "Welcome to Trivia, tonight's game will include:\n"
	if trivia.ImageRoundURL != "" {
		overviewMessage += "- An Image Round\n"
	}

	overviewMessage += fmt.Sprintf("- %d Rounds of regular trivia\n", len(trivia.Rounds))

	//TODO: should just take in the path to the audio file, no need to return the binary
	if trivia.AudioBinary != "" {
		overviewMessage += "- An Audio Round\n"
	}

	err = b.hostMessageTriviaChannel(overviewMessage)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 20)

	//Image Round
	err = b.hostMessageTriviaChannel("**Image Round:**\n" + trivia.ImageRoundTheme)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	err = b.hostMessageTriviaChannel(trivia.ImageRoundDetail)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	err = b.hostMessageTriviaChannel(trivia.ImageRoundURL)
	if err != nil {
		return err
	}
	time.Sleep(b.imageRoundSleepDelay)

	// Play Trivia
	sort.Slice(trivia.Rounds, func(i, j int) bool {
		return trivia.Rounds[i].RoundNumber < trivia.Rounds[j].RoundNumber
	})

	for _, round := range trivia.Rounds {
		roundType := strings.ToUpper(round.RoundType.Name)

		questionDelay := b.questionSleepDelay
		if roundType == lightningRoundType {
			questionDelay = time.Second * 5
		}

		err = b.hostMessageTriviaChannel(fmt.Sprintf("We will now begin **Round %d**", round.RoundNumber))
		if err != nil {
			return err
		}
		time.Sleep(time.Second)

		if roundType == listRoundType {
			err = b.hostMessageTriviaChannel(listRoundInstructions)
			if err != nil {
				return err
			}
			time.Sleep(time.Second)
		}

		if round.Theme != "" {
			err = b.hostMessageTriviaChannel(fmt.Sprintf("The theme for this round is **%s**", round.Theme))
			if err != nil {
				return err
			}
			time.Sleep(time.Second)
		}

		if round.ThemeDescription != "" {
			err = b.hostMessageTriviaChannel(round.ThemeDescription)
			if err != nil {
				return err
			}
		}

		time.Sleep(b.roundStartSleepDelay)

		sort.Slice(round.Questions, func(i, j int) bool {
			return round.Questions[i].QuestionNumber < round.Questions[j].QuestionNumber
		})

		if roundType == listRoundType {
			if len(round.Questions) != 1 {
				slog.Error("list round has more or less than one question", "num_questions", len(round.Questions))
				continue
			}
			err = b.hostMessageTriviaChannel("**Question: ** " + round.Questions[0].Question)
			if err != nil {
				return err
			}
		} else {
			for _, question := range round.Questions {
				questionMessage := fmt.Sprintf("**Question %d: ** %s", question.QuestionNumber, question.Question)
				err = b.hostMessageTriviaChannel(questionMessage)
				if err != nil {
					return err
				}
				time.Sleep(questionDelay)
			}
		}
		time.Sleep(b.roundSleepDelay)
	}

	//Mega round reminder
	err = b.hostMessageTriviaChannel(megaRoundReminder)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 10)

	// Audio Round
	// if trivia.AudioBinary != "" {
	// 	//TODO: play audio
	// }

	return b.markTriviaUsed()
}

func (b *Bot) hostMessageTriviaChannel(s string) error {
	_, err := b.triviaHost.ChannelMessageSend(
		b.triviaChannelID,
		s,
	)
	return err
}

func (b *Bot) markTriviaUsed() error {
	return b.triviaService.MarkTriviaUsed(context.Background(), b.currentTrivia.Id)
}

func (b *Bot) getNextTrivia() error {
	var err error
	b.currentTrivia, err = b.triviaService.GetNewTrivia(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving trivia from db > %w", err)
	}
	return nil
}

func (b *Bot) playAudio() error {
	if b.currentTrivia.AudioFileName == "" {
		return nil
	}

	if err := b.loadSound(); err != nil {
		//TODO send message to channel that audio isn't working now?
		return err
	}

	//TODO: would now want to spin up each audio bot and execute below code for each one

	vc, err := b.triviaHost.ChannelVoiceJoin("", "", false, false)
	if err != nil {
		return err
	}
	defer vc.Disconnect()

	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range b.audioBuffer {
		vc.OpusSend <- buff
	}

	return nil
}

func (b *Bot) loadSound() error {
	//TODO: probably want a mutex lock on this func

	cmdString := fmt.Sprintf("ffmpeg -i %s -f s16le -ar 48000 -ac 2 pipe:1 | dca > audio.dca", b.currentTrivia.AudioFileName)
	_, err := exec.Command(cmdString).Output()
	if err != nil {
		return err
	}

	file, err := os.Open("audio.dca")
	if err != nil {
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		b.audioBuffer = append(b.audioBuffer, InBuf)
	}
}
