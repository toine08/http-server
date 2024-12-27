-- +goose Up
CREATE TABLE users(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP, 
    updated_at TIMESTAMP,
    email TEXT

);

-- +goose Down
DROP TABLE users;