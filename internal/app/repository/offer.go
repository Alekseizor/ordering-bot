package repository

import (
	"context"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
)

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
