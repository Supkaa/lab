-- +goose Up
-- +goose StatementBegin
CREATE TABLE endpoints
(
    id     uuid PRIMARY KEY,
    url    TEXT NOT NULL,
    method TEXT NOT NULL
);
CREATE UNIQUE INDEX url_method_unique_idx ON endpoints (url, method);

CREATE TABLE roles
(
    name TEXT PRIMARY KEY CHECK ( LENGTH(TRIM(name)) > 0 )
);

CREATE TABLE roles_have_endpoints
(
    role_name   uuid REFERENCES roles (name),
    endpoint_id uuid REFERENCES endpoints (id)
);
CREATE UNIQUE INDEX role_name_endpoint_id_unique_idx ON roles_have_endpoints (role_name, endpoint_id);

CREATE TABLE users
(
    email    TEXT PRIMARY KEY CHECK ( LENGTH(TRIM(name)) > 0 ),
    name     TEXT NOT NULL CHECK ( LENGTH(TRIM(name)) > 0 ),
    password TEXT NOT NULL CHECK ( LENGTH(TRIM(name)) > 0 ),
    role     TEXT NOT NULL REFERENCES roles (name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

DROP INDEX role_name_endpoint_id_unique_idx;
DROP TABLE IF EXISTS roles_have_endpoints;

DROP TABLE IF EXISTS roles;

DROP INDEX url_method_unique_idx;
DROP TABLE IF EXISTS endpoints;
-- +goose StatementEnd
