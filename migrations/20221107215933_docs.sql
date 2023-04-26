-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS docs
(
    docs_id              serial PRIMARY KEY,
    docs_url             text[],
    docs_title           text[],
    images_url           text[],
    attachment           text,
    order_id                   int UNIQUE,
    chat_id                   int,
    constraint fk_customer
        FOREIGN KEY(order_id)
            REFERENCES orders(id)
                ON DELETE CASCADE,
    constraint fk_docs_conversations
        FOREIGN KEY(chat_id)
            REFERENCES conversations(id)
            ON DELETE CASCADE
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS docs;
-- +goose StatementEnd
