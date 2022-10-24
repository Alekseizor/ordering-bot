-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS offers
(
    executor_vk_id int not null,
    order_id       int not null,
    price          bigint CHECK ( price > 0 ),
    description    text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS offers;
-- +goose StatementEnd
