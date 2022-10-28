-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    id                   serial PRIMARY KEY,
    customer_vk_id       int                      NOT NULL,
    customers_comment    text,
    executor_vk_id       int,
    discipline_id        int  NOT NULL,
    date_order          timestamp with time zone,
    date_finish          timestamp with time zone,
    price                bigint CHECK ( price > 0 ),
    payout_admin         bool,
    payout_executors     bool
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
