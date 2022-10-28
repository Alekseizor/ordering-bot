package ds

import (
	"time"
)

type Order struct {
	Id               int       `db:"id"`
	CustomerVkID     int       `db:"customer_vk_id"`
	CustomersComment *string   `db:"customers_comment"`
	ExecutorVkID     *int      `db:"executor_vk_id"`
	DisciplineID     int       `db:"discipline_id"`
	DateOrder        time.Time `db:"date_order"`
	DateFinish       time.Time `db:"date_finish"`
	TimeFinish       time.Time `db:"time_finish"`
	Price            *uint64   `db:"price"`
	PayoutAdmin      *bool     `db:"payout_admin"`
	PayoutExecutors  *bool     `db:"payout_executors"`
}
