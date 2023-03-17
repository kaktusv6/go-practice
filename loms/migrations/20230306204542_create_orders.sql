-- +goose Up
-- +goose StatementBegin
create table if not exists orders(
    id serial primary key,
    user_id integer not null,
    status text not null
);

comment on table orders is 'Таблица заказов';
comment on column orders.user_id is 'ID покупателя';
comment on column orders.status is 'Статуса заказа';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders;
drop table if exists statusses;
-- +goose StatementEnd
