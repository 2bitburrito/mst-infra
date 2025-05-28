-- +goose Up
-- +goose StatementBegin

INSERT INTO licences (
  last_used_at,
  licence_type,
  user_id
) VALUES ( 
  NOW(),
  'paid',
  'e93919be-10a1-70c8-385b-c006812f9142'
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM licences
WHERE user_id='e93919be-10a1-70c8-385b-c006812f9142';
-- +goose StatementEnd
