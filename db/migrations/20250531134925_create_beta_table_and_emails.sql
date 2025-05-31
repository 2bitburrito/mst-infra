-- +goose Up
-- +goose StatementBegin
CREATE TABLE beta_licences (
  email TEXT
);
INSERT INTO beta_licences (email)
VALUES
    ('palmerhap@gmail.com'),
    ('davidrossaudio@gmail.com'),
    ('cutsnake@netspace.net.au'),
    ('Ale.cesana88@gmail.com'),
    ('ryan@deadonsound.com'),
    ('paul_isaacs@sounddevices.com'),
    ('josh@camelcitysound.net'),
    ('lukearmstrong@me.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS beta_licences;
-- +goose StatementEnd
