package ds

type Offer struct {
	OfferID      int `db:"offer_id"`
	ExecutorVKID int `db:"executor_vk_id"`
	OrderID      int `db:"order_id"`
	Price        int `db:"price"`
}
