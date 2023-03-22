-- +goose Up
-- +goose StatementBegin
alter table if exists orders
    add column if not exists created_at timestamp;

update orders set created_at = now() where created_at is null;

alter table if exists orders
    alter column created_at set not null;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists orders
    drop column if exists created_at;
-- +goose StatementEnd
