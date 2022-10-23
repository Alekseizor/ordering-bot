package ds

import "time"

type Order struct {
	Id                  int
	CustomerVkID        int        `db:"customer_vk_id"`
	DisciplineID        int        `db:"discipline_id"`
	DisciplineName      string     `db:"discipline_id"`
	ExecutorVkID        *int       `db:"executor_vk_id"`
	CreatedAt           time.Time  `db:"created_at"`
	DeadlineAt          time.Time  `db:"deadline_at"`
	DoneAt              *time.Time `db:"done_at"`
	Paid                bool
	Price               *uint64
	PercentageDeduction int    `db:"percentage_deduction"`
	Comment             string `db:"customers_comment"`
}
