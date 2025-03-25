package usecase

import (
	"fmt"

	"github.com/brkss/dextrace/internal/domain"
)

type SibionicUseCase struct {
	authRepo    domain.AuthRepository
	glucoseRepo domain.GlucoseRepository
}

func NewSibionicUseCase(authRepo domain.AuthRepository, glucoseRepo domain.GlucoseRepository) *SibionicUseCase {
	return &SibionicUseCase{
		authRepo:    authRepo,
		glucoseRepo: glucoseRepo,
	}
}

func (uc *SibionicUseCase) GetGlucoseData(user domain.User, userID string) (*domain.GlucoseDataResponse, error) {
	token, err := uc.authRepo.Login(user)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	data, err := uc.glucoseRepo.GetData(token, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get glucose data: %w", err)
	}

	return data, nil
}