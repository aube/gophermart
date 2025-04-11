-- +goose Up
-- +goose StatementBegin

CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders (
    id bigint not null primary key,
    order_id bigint not null,
    user_id serial not null,
    loyalty_points bigint not null check (loyalty_points >= 0) default 0,
    status order_status NOT NULL DEFAULT 'NEW'
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean not null default false
);

create INDEX order_id on orders (order_id);
create INDEX user_id on orders (user_id);


CREATE TRIGGER orders_updated_at_trigger
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER orders_updated_at_trigger ON orders;

drop INDEX order_id;
drop INDEX user_id;

DROP TABLE orders;

-- +goose StatementEnd
