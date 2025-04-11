-- +goose Up
-- +goose StatementBegin

CREATE TABLE withdravals (
    id serial not null primary key,
    user_id int not null,
    amount bigint not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create INDEX user_id on withdravals (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop INDEX user_id;

DROP TABLE withdravals;

-- +goose StatementEnd
