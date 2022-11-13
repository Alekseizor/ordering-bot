-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    id                   serial PRIMARY KEY,
    customer_vk_id       int                      NOT NULL,
    customers_comment    text,
    executor_vk_id       int,
    type_order           text NOT NULL,
    discipline_id        int,
    date_order           timestamp with time zone,
    date_finish          timestamp with time zone,
    price                bigint CHECK ( price > 0 ),
    percent_executor     int,
    verification_executor         bool,
    verification_customer        bool,
    order_task           text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
