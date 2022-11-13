-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS executors
(
    id                   serial PRIMARY KEY,
    vk_id                int CHECK ( vk_id >= 0 ) UNIQUE,
    disciplines_id       int[],
    percent_executor  int CHECK ( percent_executor >= 0 and percent_executor<=100),
    rating               float CHECK ( rating >= 0 ) DEFAULT(0),
    amount_rating        int CHECK ( amount_rating >= 0 ) DEFAULT(0),
    profit               float CHECK ( profit >= 0 )  DEFAULT(0),
    amount_orders        int CHECK ( amount_orders >= 0 ) DEFAULT(0),
    requisite      text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS executors;
-- +goose StatementEnd
