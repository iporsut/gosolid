-- +goose Up
-- +goose StatementBegin
CREATE TABLE post (
    id INT NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE post;
-- +goose StatementEnd
