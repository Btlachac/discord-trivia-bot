package worker

import (
	"fmt"
	"go-trivia-api/internal/db"
	"log/slog"
	"sort"
	"strings"
	"sync"
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

// https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go seems useful

type triviaService interface {
	GetNewTrivia() (db.Trivia, error)
	MarkTriviaUsed(triviaId int64) error
}

// TODO: add mutex lock around trivia in progress
// TODO: do validation on channel IDs in main
type bot struct {
	triviaHost    *discordgo.Session
	triviaService triviaService

	triviaChannelID  string
	commandChannelID string

	currentTrivia db.Trivia
	triviaLock    sync.Mutex

	imageRoundSleepDelay time.Duration
	roundSleepDelay      time.Duration
	questionSleepDelay   time.Duration
	roundStartSleepDelay time.Duration
}

func (b *bot) Run() {
	b.triviaHost.AddHandler(b.messageCreateHandler)
}

// general func to listen for messages
func (b *bot) messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		//TODO
	case "!audio":
		//TODO
	case "!image":
		//TODO
	}
}

// TODO: think about if we want to return for any error? If we miss 1 message it might be better to just log it? Perhaps create a budget,
// if more than 5 messages fail we die
func (b *bot) startTrivia() error {
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
	b.currentTrivia, err = b.triviaService.GetNewTrivia()
	if err != nil {
		return fmt.Errorf("error retrieving trivia from db > %w", err)
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

func (b *bot) hostMessageTriviaChannel(s string) error {
	_, err := b.triviaHost.ChannelMessageSend(
		b.triviaChannelID,
		s,
	)
	return err
}

func (b *bot) markTriviaUsed() error {
	return b.triviaService.MarkTriviaUsed(b.currentTrivia.Id)
}
