-- +goose Up

CREATE TABLE USERS (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  hashed_password CHAR(60) NOT NULL
);

CREATE TABLE FRANCHISES (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,

  user_id INTEGER REFERENCES USERS(id) NOT NULL
);

CREATE TABLE GAMES (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,

  franchise_id INTEGER REFERENCES FRANCHISES(id),
  
  status VARCHAR(255) DEFAULT 'To Do',
  current_progress INTEGER DEFAULT 0,
  final_progress INTEGER DEFAULT 100,

  user_id INTEGER REFERENCES USERS(id) NOT NULL
);

CREATE TABLE USER_TOKENS (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES USERS(id),
  token VARCHAR(400)
);

ALTER TABLE games ADD CONSTRAINT games_uc_name_user_id UNIQUE (name, user_id);
ALTER TABLE games ADD CONSTRAINT games_uc_name_user_id_franchise UNIQUE (name, user_id, franchise_id);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
ALTER TABLE users ADD CONSTRAINT users_uc_username UNIQUE (username);

ALTER TABLE franchises ADD CONSTRAINT franchises_uc_name_user_id UNIQUE (name, user_id);

-- +goose Down
DROP TABLE USER_TOKENS;
DROP TYPE GAME_STATUS;
DROP TABLE USERS;
DROP TABLE GAMES;
DROP TABLE FRANCHISES;

