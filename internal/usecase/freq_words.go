package usecase

import (
	"context"
	"fmt"
	"npp_doslab/internal/entity"
	"npp_doslab/internal/usecase/repo"
)

type FrequentlyWordsUseCase struct {
	repo *repo.FreqWordsRepo
}

func New(r *repo.FreqWordsRepo) *FrequentlyWordsUseCase {
	return &FrequentlyWordsUseCase{
		repo: r,
	}
}

func (uc *FrequentlyWordsUseCase) GetByPostId(postId int, ctx context.Context) ([]entity.Comm, error) {
	words, err := uc.repo.GetWordsStatisticsByPostId(postId, ctx)
	if err != nil {
		return nil, fmt.Errorf("FrequentlyWordsUseCase - Scan - s.repo.Scan: %w", err)
	}

	return words, nil
}
