package repository

import (
	"context"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
)

func (r *Repository) NewOffer(ctx context.Context, offer ds.Offer) error {
	var err error

	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	_, err = r.db.ExecContext(ctx, "INSERT INTO offers "+
		"(executor_vk_id, order_id, price)"+
		"VALUES ($1, $2, $3)", offer.ExecutorVkID, offer.OrderID, offer.Price)
	if err != nil {
		return err
	}

	return nil
}

// AvailableOrders показывает какие пользователю доступны заказы
func (r *Repository) AvailableOrders(ctx context.Context, vkID int) ([]*ds.Order, error) {
	var err error
	var orders []*ds.Order

	// start ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	q := r.sq.Select("*").From("orders").Where(
		"orders.discipline_id = ANY (ARRAY (SELECT disciplines_id FROM executors WHERE executors.vk_id = $1))", vkID)

	qSelect, argsSelect, err := q.
		OrderBy("id").
		ToSql()

	err = r.db.SelectContext(ctx, &orders, qSelect, argsSelect...)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
