CREATE TABLE IF NOT EXISTS GAMES (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS USERS (
  id SERIAL PRIMARY KEY, 
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  hashed_password CHAR(60) NOT NULL
);

CREATE TYPE game_status AS ENUM ('To Do', 'In Progress', 'Done');

CREATE TABLE IF NOT EXISTS USER_GAMES (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES USERS(id),
  game_id INTEGER REFERENCES GAMES(id),
  status game_status DEFAULT 'To Do'
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
ALTER TABLE users ADD CONSTRAINT users_uc_username UNIQUE (username);

ALTER TABLE user_games ADD CONSTRAINT user_games_unique_user_game (user_id, game_id);
