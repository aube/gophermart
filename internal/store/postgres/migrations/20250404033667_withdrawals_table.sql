-- +goose Up
-- +goose StatementBegin

CREATE TABLE withdrawals (
    id serial not null primary key,
    user_id int not null,
    accrual bigint not null check (accrual >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create INDEX withdrawals_user_id on withdrawals (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop INDEX withdrawals_user_id;

DROP TABLE withdrawals;

-- +goose StatementEnd
