-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS executors
(
    vk_id          int CHECK ( vk_id >= 0 ) UNIQUE,
    disciplines_id int[],
    proportion     int CHECK ( proportion >= 0 and proportion<=100),
    rating         float NOT NULL DEFAULT(5)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS executors;
-- +goose StatementEnd
