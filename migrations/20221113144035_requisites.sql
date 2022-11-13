-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS requisites
(
    requisites_id     serial PRIMARY KEY,
    requisites          text
);
-- +goose StatementEnd
INSERT INTO requisites (requisites)
VALUES ('Сбербанк:\n5469 3800 6905 9201\nТинькофф:\n5536 9139 6190 8197\nВладимир Михайлович К.');
-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS requisites;
-- +goose StatementEnd
