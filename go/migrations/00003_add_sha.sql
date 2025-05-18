-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents ADD COLUMN sha VARCHAR(64);
ALTER TABLE memories ADD COLUMN sha VARCHAR(64);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN sha;
ALTER TABLE memories DROP COLUMN sha;
-- +goose StatementEnd
