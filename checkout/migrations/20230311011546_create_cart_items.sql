-- +goose Up
-- +goose StatementBegin
create table if not exists cart_items(
    user_id int not null,
    sku int not null,
    count int not null
);

comment on table cart_items is 'Таблица с товарами в корзинах пользователей';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists cart_items;
-- +goose StatementEnd
