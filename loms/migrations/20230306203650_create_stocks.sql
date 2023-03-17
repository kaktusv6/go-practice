-- +goose Up
-- +goose StatementBegin
create table if not exists stocks(
    sku integer not null,
    warehouse_id integer not null,
    count integer not null
);

comment on table stocks is 'Таблица остатков товаров';
comment on column stocks.sku is 'Артикул товара';
comment on column stocks.warehouse_id is 'ID склада';
comment on column stocks.count is 'Кол-во товара';

insert into stocks (sku, warehouse_id, count)
values
    (4996014, 1, 100),
    (4996014, 2, 1),
    (5097510, 2, 1),
    (5415913, 3, 1),
    (5647362, 4, 2),
    (6147564, 5, 5),
    (6245113, 1, 15),
    (6245113, 1, 56),
    (6966051, 2, 20),
    (6967749, 3, 20),
    (7277168, 4, 20),
    (7277168, 2, 2),
    (7895903, 5, 60),
    (8748527, 1, 45),
    (8748527, 5, 5),
    (18247421, 2, 65),
    (19045918, 3, 58),
    (19065113, 1, 1),
    (19065113, 2, 2),
    (19065113, 3, 3),
    (19065113, 4, 4),
    (19065113, 5, 5),
    (19366373, 3, 76),
    (19717106, 5, 19),
    (20367334, 2, 53),
    (21435785, 4, 73),
    (24015704, 3, 41),
    (24166967, 1, 58),
    (24167411, 5, 2),
    (24167617, 2, 5),
    (24167788, 4, 6),
    (24168225, 3, 4),
    (24418527, 3, 7),
    (24438546, 2, 9),
    (24438546, 5, 90),
    (24768684, 3, 8),
    (24808287, 1, 3),
    (24808287, 2, 30),
    (25475334, 4, 21),
    (26176267, 5, 11),
    (26176267, 1, 3),
    (26176267, 3, 34),
    (26267884, 2, 4)
;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists stocks;
-- +goose StatementEnd
