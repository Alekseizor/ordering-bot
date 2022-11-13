-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS newsletter
(
    newsletter_id serial PRIMARY KEY,
    docs_url      text[],
    docs_title    text[],
    images_url    text[],
    attachment    text,
    text_message  text,
    peer_ids      int[]
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS newsletter;
-- +goose StatementEnd
