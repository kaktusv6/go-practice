-- +goose Up
-- +goose StatementBegin
create table if not exists orders_items_stocks(
    order_id int not null,
    sku int not null,
    count int not null,
    warehouse_id int not null
);

comment on table orders_items_stocks is 'Таблица с резервированием товаров со склада для заказов';
comment on column orders_items_stocks.order_id is 'ID заказа (orders.id)';
comment on column orders_items_stocks.sku is 'Артикул товара';
comment on column orders_items_stocks.count is 'Кол-во товаров';
comment on column orders_items_stocks.warehouse_id is 'ID склада с которого резеврируется товар';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists order_items_stocks;
-- +goose StatementEnd
