-- +goose Up
-- +goose StatementBegin
create table if not exists orders_items(
    order_id int not null,
    sku int not null,
    warehouse_id int not null,
    count int not null
);

comment on table orders_items is 'Таблица зарезервированных товаров';
comment on column orders_items.order_id is 'ID аказа (orders)';
comment on column orders_items.sku is 'Артикул заказа';
comment on column orders_items.warehouse_id is 'ID склада на котором лежит товар';
comment on column  orders_items.count is 'Кол-во заразервированного заказа';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders_items;
-- +goose StatementEnd
