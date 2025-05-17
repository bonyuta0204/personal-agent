-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents ADD COLUMN modified_at TIMESTAMP;
ALTER TABLE memories ADD COLUMN modified_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN modified_at;
ALTER TABLE memories DROP COLUMN modified_at;
-- +goose StatementEnd
