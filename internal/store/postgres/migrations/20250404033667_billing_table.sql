-- +goose Up
-- +goose StatementBegin

CREATE TABLE billing (
    id serial not null primary key,
    order_id bigint not null,
    user_id int not null,
    amount bigint not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create INDEX billing_user_id on billing (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop INDEX billing_user_id;

DROP TABLE billing;

-- +goose StatementEnd
