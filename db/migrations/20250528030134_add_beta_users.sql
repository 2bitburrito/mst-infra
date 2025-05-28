-- +goose Up
-- +goose StatementBegin
    INSERT INTO users (email, has_license, number_of_licenses, subscribed_to_emails, full_name, id)
    VALUES
        ('palmerhap@gmail.com', true, 1, true, 'Hugh Palmer', gen_random_uuid()),
        ('davidrossaudio@gmail.com', true, 1, true, 'David Ross', gen_random_uuid()),
        ('cutsnake@netspace.net.au', true, 1, true, 'Terry Chadwick', gen_random_uuid()),
        ('Ale.cesana88@gmail.com', true, 1, true, 'Ale Cesana', gen_random_uuid()),
        ('ryan@deadonsound.com', true, 1, true, 'Ryan Granger', gen_random_uuid()),
        ('paul_isaacs@sounddevices.com', true, 1, true, 'Paul Isaacs', gen_random_uuid()),
        ('josh@camelcitysound.net', true, 1, true, 'Josh Tucker', gen_random_uuid()),
        ('lukearmstrong@me.com', true, 1, true, 'Luke Armstrong', gen_random_uuid());

    INSERT INTO licences (user_id, licence_type, expiry, last_used_at)
    SELECT id, 'beta', NOW() + INTERVAL '90 days', NOW()
    FROM users
    WHERE email IN (
        'palmerhap@gmail.com',
        'davidrossaudio@gmail.com',
        'cutsnake@netspace.net.au',
        'Ale.cesana88@gmail.com',
        'ryan@deadonsound.com',
        'paul_isaacs@sounddevices.com',
        'josh@camelcitysound.net',
        'lukearmstrong@me.com'
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    DELETE FROM licences
    WHERE user_id IN (
        SELECT id
        FROM users
        WHERE email IN (
            'palmerhap@gmail.com',
            'davidrossaudio@gmail.com',
            'cutsnake@netspace.net.au',
            'Ale.cesana88@gmail.com',
            'ryan@deadonsound.com',
            'paul_isaacs@sounddevices.com',
            'josh@camelcitysound.net',
            'lukearmstrong@me.com'
        )
    );

    DELETE FROM users
    WHERE email IN (
        'palmerhap@gmail.com',
        'davidrossaudio@gmail.com',
        'cutsnake@netspace.net.au',
        'Ale.cesana88@gmail.com',
        'ryan@deadonsound.com',
        'paul_isaacs@sounddevices.com',
        'josh@camelcitysound.net',
        'lukearmstrong@me.com'
    );
-- +goose StatementEnd
