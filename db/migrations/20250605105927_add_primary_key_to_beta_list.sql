-- +goose Up
-- +goose StatementBegin
ALTER TABLE beta_licences
   ADD PRIMARY KEY(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE beta_licences
   DROP CONSTRAINT beta_licences_pkey;
-- +goose StatementEnd
