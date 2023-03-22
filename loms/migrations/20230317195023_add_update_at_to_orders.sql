-- +goose Up
-- +goose StatementBegin
alter table if exists orders
    add column if not exists updated_at timestamp;

update orders set updated_at = created_at where updated_at is null;

alter table if exists orders
    alter column updated_at set not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists orders
    drop column if exists updated_at;
-- +goose StatementEnd
