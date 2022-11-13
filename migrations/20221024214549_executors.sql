-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS executors
(
    vk_id          int CHECK ( vk_id >= 0 ) UNIQUE,
    disciplines_id int[],
    precent_executor  int CHECK ( precent_executor >= 0 and precent_executor<=100),
    rating         float CHECK ( rating >= 0 ) DEFAULT(0),
    amount_rating  int CHECK ( amount_rating >= 0 ) DEFAULT(0),
    profit         float CHECK ( profit >= 0 )  DEFAULT(0),
    amount_orders  int CHECK ( amount_orders >= 0 ) DEFAULT(0),
    requisite      text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS executors;
-- +goose StatementEnd
