package ds

import (
	"time"
)

type Order struct {
	Id                   int       `db:"id"`
	CustomerVkID         uint      `db:"customer_vk_id"`
	CustomersComment     *string   `db:"customers_comment"`
	ExecutorVkID         *uint     `db:"executor_vk_id"`
	TypeOrder            string    `db:"type_order"`
	DisciplineID         int       `db:"discipline_id"`
	DateOrder            time.Time `db:"date_order"`
	DateFinish           time.Time `db:"date_finish"`
	Price                *uint64   `db:"price"`
	PercentExecutor      *uint     `db:"percent_executor"`
	VerificationExecutor *bool     `db:"verification_executor"`
	VerificationCustomer *bool     `db:"verification_customer"`
	OrderTask            *string   `db:"order_task"`
}
