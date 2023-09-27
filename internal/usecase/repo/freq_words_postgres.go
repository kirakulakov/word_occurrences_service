package repo

import (
	"context"
	"fmt"
	"npp_doslab/internal/entity"
	"npp_doslab/pkg/logger"
	"npp_doslab/pkg/postgres"
)

const (
	_defaultEntityCap     = 100
	_noRowsInResultSetExc = "no rows in result set"
)

type FreqWordsRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *FreqWordsRepo {
	return &FreqWordsRepo{pg}
}

func (r *FreqWordsRepo) GetWords(ctx context.Context) ([]entity.Comm, error) {
	sql, _, err := r.Builder.
		Select("post_id, word", "count").
		From("word").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("FreqWordsRepo - GetWords - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("FreqWordsRepo - GetWords - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.Comm, 0, _defaultEntityCap)

	for rows.Next() {
		e := entity.Comm{}

		err = rows.Scan(&e.PostId, &e.Word, &e.Count)
		if err != nil {
			return nil, fmt.Errorf("FreqWordsRepo - GetWords - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (r *FreqWordsRepo) _getCurrentCountOccurrencesOfWordInDB(word string, ctx context.Context, l *logger.Logger) int {
	sql, args, err := r.Builder.
		Select("count").
		From("word").
		Where("word = ?", word).
		ToSql()
	if err != nil {
		l.Fatal(fmt.Errorf("FreqWordsRepo - InsertCommentData - r.Builder: %w", err))
	}

	rows := r.Pool.QueryRow(ctx, sql, args...)

	c := entity.Count{}

	err = rows.Scan(&c.Count)
	if err != nil {
		if err.Error() == _noRowsInResultSetExc {
			return 0
		} else {
			l.Fatal(fmt.Errorf("FreqWordsRepo - InsertCommentData - rows.Scan: %w", err))
		}
	}

	return c.Count

}

func (r *FreqWordsRepo) _addNew(
	postId int, word string, countOccurrences int, ctx context.Context, l *logger.Logger) {

	sql, args, err := r.Builder.
		Insert("word").
		Columns("post_id", "word", "count").
		Values(postId, word, countOccurrences).
		ToSql()
	if err != nil {
		l.Fatal(fmt.Errorf("FreqWordsRepo - InsertCommentData - r.Builder: %w", err))
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		l.Fatal(fmt.Errorf("FreqWordsRepo - InsertCommentData - r.Pool.Query: %w", err))
	}

	rows.Close()
}

func (r *FreqWordsRepo) _updateExisting(
	currentCount int, word string, countOccurrences int, ctx context.Context, l *logger.Logger) {

	new_count := currentCount + countOccurrences
	sql, args, err := r.Builder.Update("word").Set("count", new_count).Where("word = ?", word).ToSql()

	if err != nil {
		l.Fatal(fmt.Errorf("FreqWordsRepo - InsertCommentData - r.Builder: %w", err))
	}
	_, err = r.Pool.Exec(ctx, sql, args...)

	if err != nil {
		l.Fatal(fmt.Errorf("FreqWordsRepo - UpdateData - r.Pool.Exec: %w", err))
	}
}

func (r *FreqWordsRepo) UpdateStatisticInDB(
	wordsByCountOccurrences map[string]int, postId int, ctx context.Context, l *logger.Logger) {

	for word, count_occurences := range wordsByCountOccurrences {

		currentCount := r._getCurrentCountOccurrencesOfWordInDB(word, ctx, l)

		if currentCount == 0 {
			r._addNew(postId, word, count_occurences, ctx, l)
		} else {
			r._updateExisting(currentCount, word, count_occurences, ctx, l)

		}

	}

}
