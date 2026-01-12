CREATE TABLE users(
    'id' INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    'name' VARCHAR(255) NOT NULL,
    'email' VARCHAR(255) NOT NULL UNIQUE,
    'password' VARCHAR(255) NOT NULL,
    'is_admin' BOOLEAN NOT NULL DEFAULT FALSE,
    'created_at' datetime DEFAULT (now()),
    'updated_at' datetime,
)