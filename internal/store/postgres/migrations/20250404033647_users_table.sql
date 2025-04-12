-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id serial not null primary key,
    email varchar not null unique,
    encrypted_password varchar not null,
    accrual bigint not null check (accrual >= 0) default 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean not null default false
);

create INDEX users_email on users (email);

create INDEX users_encrypted_password on users (encrypted_password);


CREATE TRIGGER users_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER users_updated_at_trigger ON users;

drop INDEX users_encrypted_password;
drop INDEX users_email;

DROP TABLE users;

-- +goose StatementEnd
