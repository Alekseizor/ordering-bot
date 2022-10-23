package repository

import (
	"context"
	"fmt"
	"github.com/Alekseizor/ordering-bot/internal/app/config"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

func (r *Repository) IsExecutor(ctx context.Context, vkID int, disciplineID ...int) (bool, error) {
	var err error
	var exists []int
	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	//получили число записей в таблице executors
	q := r.sq.Select("count(*)").From("executors")

	if disciplineID != nil {
		q = q.Where("? = ANY (disciplines_id)", disciplineID[0])
	}

	qSelect, argsSelect, err := q.Where(sq.Eq{"vk_id": vkID}).
		Limit(1).
		ToSql()

	err = r.db.SelectContext(ctx, &exists, qSelect, argsSelect...)
	if err != nil {
		return false, err
	}

	if exists[0] == 1 {
		return true, nil
	}

	return false, nil
}

func (r *Repository) GetPotentialExecutors(ctx context.Context, disciplineID int) ([]*int, error) {
	var err error
	var executors_ids []*int

	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	q := r.sq.Select().From("executors")
	qSelect, argsSelect, err := q.Columns("vk_id").
		OrderBy("rating").
		Where("? = ANY (disciplines_id)", disciplineID).
		ToSql()

	err = r.db.SelectContext(ctx, &executors_ids, qSelect, argsSelect...)
	if err != nil {
		return nil, err
	}

	return executors_ids, nil
}

func (r *Repository) IsAdmin(ctx context.Context, vkID int) bool {
	return config.FromContext(ctx).AdminID == vkID
}

func (r *Repository) RemoveExecutor(ctx context.Context, vkID int) error {
	qDelete, _, err := r.sq.Delete("executors").Where("vk_id = ?", vkID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, qDelete, vkID)

	return err
}

func (r *Repository) CreateExecutor(ctx context.Context, e ds.Executor) error {
	var err error

	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	_, err = r.db.ExecContext(ctx, "INSERT INTO executors (vk_id, disciplines_id, rating) "+
		"VALUES ($1, $2, $3)", e.VkID, pq.Array(e.Disciplines), 5)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetExecutorsDisciplines(ctx context.Context, vkID int) ([]*ds.Discipline, error) {
	var err error
	var disciplines []*ds.Discipline

	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	q := r.sq.Select("id, name, percent").
		From(fmt.Sprintf("unnest(ARRAY(SELECT disciplines_id FROM executors WHERE vk_id = %d)) did LEFT JOIN disciplines t on t.id = did", vkID))
	qSelect, argsSelect, err := q.
		OrderBy("id").
		ToSql()

	err = r.db.SelectContext(ctx, &disciplines, qSelect, argsSelect...)
	if err != nil {
		return nil, err
	}

	return disciplines, nil
}
