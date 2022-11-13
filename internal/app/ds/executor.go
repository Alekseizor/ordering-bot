package ds

type Executor struct {
	Id              int     `db:"id"`
	VkID            int     `db:"vk_id"`
	DisciplinesID   []int   `db:"disciplines_id"`
	PercentExecutor int     `db:"percent_executor"`
	Rating          float64 `db:"rating"`
	AmountRating    int     `db:"amount_rating"`
	Profit          float64 `db:"profit"`
	AmountOrders    int     `db:"amount_orders"`
	Requisite       *string `db:"requisite"`
}
