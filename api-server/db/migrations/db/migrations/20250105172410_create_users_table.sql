-- migrate:up
CREATE TABLE users (
    uuid UUID,
    identities TEXT,
    credits INT,
    last_request INT,
    is_premium BOOLEAN,
    PRIMARY KEY (uuid)
);

-- migrate:down

