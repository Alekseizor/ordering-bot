-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    id                   serial PRIMARY KEY,
    customer_vk_id       int                      NOT NULL,
    customers_comment    text,
    executor_vk_id       int,
    type_order           text NOT NULL DEFAULT 'Домашнее задание',
    discipline_id        int DEFAULT(0),
    date_order           timestamp with time zone DEFAULT now(),
    date_finish          timestamp with time zone DEFAULT now(),
    price                bigint CHECK ( price > 0 ),
    percent_executor     int,
    verification_executor         bool,
    verification_customer        bool,
    order_task           text DEFAULT 'Домашнее задание'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
