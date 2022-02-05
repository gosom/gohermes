CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  identifier varchar(50) NOT NULL UNIQUE
);
