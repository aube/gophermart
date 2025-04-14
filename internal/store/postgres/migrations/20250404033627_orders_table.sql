-- +goose Up
-- +goose StatementBegin

CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders (
    id bigint not null primary key,
    user_id int not null,
    accrual bigint not null check (accrual >= 0) default 0,
    status order_status NOT NULL DEFAULT 'NEW',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean not null default false
);

create INDEX orders_user_id on orders (user_id);


CREATE TRIGGER orders_updated_at_trigger
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER orders_updated_at_trigger ON orders;

drop INDEX orders_user_id;

DROP TABLE orders;

drop TYPE order_status;

-- +goose StatementEnd
