-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id      SERIAL,
    name    TEXT NOT NULL ,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
