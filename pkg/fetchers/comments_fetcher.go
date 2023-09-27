package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"npp_doslab/config"
	"npp_doslab/internal/entity"
	"npp_doslab/internal/usecase/repo"
	"npp_doslab/pkg/logger"
	"npp_doslab/pkg/postgres"
	"strings"
	"time"
)

type FetchData struct {
	data interface{}
}

type Subscription interface {
	Updates() <-chan []entity.Comment
}

type sub struct {
	fetcher Fetcher
	updates chan []entity.Comment
}

func NewSubscription(ctx context.Context, fetcher Fetcher, freq int, l *logger.Logger) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan []entity.Comment),
	}
	go s.serve(ctx, freq, l)
	return s
}

type Fetcher interface {
	Fetch(int) ([]entity.Comment, error)
}

func (s *sub) Updates() <-chan []entity.Comment {
	return s.updates
}

func NewFetcher(uri string, l *logger.Logger) Fetcher {
	f := &fetcher{
		uri: uri,
		l:   l,
	}
	return f
}

type fetcher struct {
	uri string
	l   *logger.Logger
}

func (f *fetcher) Fetch(postId int) ([]entity.Comment, error) {
	url := fmt.Sprintf(f.uri, postId)

	resp, err := http.Get(url)
	if err != nil {
		return []entity.Comment{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []entity.Comment{}, err
	}

	var comments []entity.Comment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return []entity.Comment{}, err
	}

	return comments, nil
}

func (s *sub) serve(ctx context.Context, checkFrequency int, l *logger.Logger) {
	clock := time.NewTicker(time.Duration(checkFrequency) * time.Second)
	type fetchResult struct {
		fetched []entity.Comment
		err     error
	}
	fetchDone := make(chan fetchResult, 1)

	for {
		select {
		case <-clock.C:

			for i := 0; i < 100; i++ { // Хардкод, но можно добавить дополнительный шаг получения id's всех существующих постов, и по ним проитерироваться
				go func(i int) {
					fetched, err := s.fetcher.Fetch(i)
					fetchDone <- fetchResult{fetched, err}
				}(i)
			}

		case result := <-fetchDone:
			fetched := result.fetched
			if result.err != nil {
				l.Fatal(fmt.Errorf("Fetch error: %v \n Waiting the next iteration", result.err.Error()))
				break
			}
			s.updates <- fetched
		case <-ctx.Done():
			return
		}
	}
}

func groupWordsByCountOccurrences(words []string) map[string]int {
	_wordsByCountOccurrences := make(map[string]int)
	for _, word := range words {
		_wordsByCountOccurrences[word]++
	}

	return _wordsByCountOccurrences
}

func RunFetching(freq int, cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	repo := repo.New(pg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	ctx, cancel := context.WithCancel(context.Background())
	subscription := NewSubscription(ctx, NewFetcher("https://jsonplaceholder.typicode.com/posts/%v/comments", l), freq, l)

	for comments := range subscription.Updates() {
		for i := 0; i < len(comments); i++ {

			commentPostId := comments[i].PostID
			commentBody := comments[i].Body

			wordsFromCommentBody := strings.Split(commentBody, " ")

			repo.UpdateStatisticInDB(
				groupWordsByCountOccurrences(wordsFromCommentBody),
				commentPostId,
				ctx, l)
		}

	}

	defer cancel()
}
