package ds

type Conversation struct {
	Id                      int `db:"id"`
	OrderID                 int `db:"order_id"`
	CustomersConversationID int `db:"customers_conversation_id"`
	ExecutorsConversationID int `db:"executors_conversation_id"`
}
