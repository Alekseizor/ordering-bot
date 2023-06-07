-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS conversations
(
    id                   serial PRIMARY KEY,
    order_id       int                      NOT NULL,
    customers_conversation_id    int,
    executors_conversation_id       int
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS conversations;
-- +goose StatementEnd
