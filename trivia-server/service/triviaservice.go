package service

import (
	"go-trivia-api/model"
	"go-trivia-api/repository"
)

type TriviaService struct {
	triviaRepository *repository.TriviaRepository
}

func NewTriviaService(triviaRepository *repository.TriviaRepository) *TriviaService {
	return &TriviaService{
		triviaRepository: triviaRepository,
	}
}

func (service *TriviaService) GetNewTrivia() model.Trivia {
	return service.triviaRepository.GetNewTrivia()
}

func (service *TriviaService) AddTrivia(newTrivia model.Trivia) {
	service.triviaRepository.AddTrivia(newTrivia)
}

func (service *TriviaService) MarkTriviaUsed(triviaId int64) {
	service.triviaRepository.MarkTriviaUsed(triviaId)
}

func (service *TriviaService) GetTriviaList() []model.TriviaInfo {
	return service.triviaRepository.GetTriviaList()
}
