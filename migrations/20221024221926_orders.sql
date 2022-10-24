-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    id                   uuid,
    customer_vk_id       int                      NOT NULL,
    customers_comment    text NOT NULL DEFAULT(''),
    executor_vk_id       int,
    discipline_id        int  NOT NULL,
    date_finish          timestamp with time zone NOT NULL,
    time_finish          timestamp with time zone NOT NULL,
    price                bigint CHECK ( price > 0 ),
    payout_admin         bool,
    payout_executors     bool,
    status               int                      NOT NULL DEFAULT (1)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
