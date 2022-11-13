package ds

type Newsletter struct {
	NewsletterID int      `db:"newsletter_id"`
	DocsUrl      []string `db:"docs_url"`
	DocsTitle    []string `db:"docs_title"`
	ImagesUrl    []string `db:"images_url"`
	Attachment   *string  `db:"attachment"`
	TextMessage  string   `db:"text_message"`
	PeerIDs      []int    `db:"peer_ids"`
}
