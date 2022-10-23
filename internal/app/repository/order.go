package repository

import (
	"context"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CreateNewOrder(ctx context.Context, o ds.Order) error {
	var err error

	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	_, err = r.db.ExecContext(ctx, "INSERT INTO orders "+
		"(paid, discipline_id, customer_vk_id, customers_comment, percentage_deduction, deadline_at) "+
		"VALUES ($1, $2, $3, $4, $5, $6)", false, o.DisciplineID, o.CustomerVkID, o.Comment, 0, o.DeadlineAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserOrders(ctx context.Context, userVKID int, status ...ds.OrderStatus) ([]*ds.Order, error) {
	var err error
	var orders []*ds.Order

	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	q := r.sq.Select("*").From("orders")

	if status != nil {
		q = q.Where(sq.Eq{"status": status[0]})
	}

	qSelect, argsSelect, err := q.
		Where(sq.Eq{"customer_vk_id": userVKID}).
		OrderBy("id").
		ToSql()

	err = r.db.SelectContext(ctx, &orders, qSelect, argsSelect...)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) GetAllOrders(ctx context.Context) ([]*ds.Order, error) {
	var err error
	var orders []*ds.Order

	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	q := r.sq.Select("*").From("orders")

	qSelect, argsSelect, err := q.
		OrderBy("id").
		ToSql()

	err = r.db.SelectContext(ctx, &orders, qSelect, argsSelect...)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) GetOrderByID(ctx context.Context, id int) (ds.Order, error) {
	var err error
	var order []*ds.Order

	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	q := r.sq.Select("*").From("orders")

	qSelect, argsSelect, err := q.
		OrderBy("id").
		Limit(1).
		ToSql()

	err = r.db.SelectContext(ctx, &order, qSelect, argsSelect...)
	if err != nil {
		return ds.Order{}, err
	}

	return *order[0], nil
}
