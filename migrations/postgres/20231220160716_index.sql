-- +goose Up
-- +goose StatementBegin
CREATE INDEX surname_idx ON Employees(surname);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX surname_idx;
-- +goose StatementEnd
