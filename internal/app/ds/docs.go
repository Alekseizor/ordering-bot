package ds

type Docs struct {
	DocsID     int      `db:"docs_id"`
	DocsUrl    []string `db:"docs_url"`
	DocsTitle  []string `db:"docs_title"`
	ImagesUrl  []string `db:"images_url"`
	Attachment *string  `db:"attachment"`
	OrderID    int      `db:"order_id"`
	//ChatID     int      `db:"chat_id"`
}
