
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS Users (
  id uuid DEFAULT uuid_generate_v4 (),
  username TEXT NOT NULL,
  password_hash TEXT NOT NULL,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS Expenses (
  id serial PRIMARY KEY,
  title TEXT NOT NULL,
  amount decimal NOT NULL,
  note TEXT,
	tags TEXT[],
  user_id uuid,

  CONSTRAINT fk_user
  FOREIGN KEY (user_id) REFERENCES Users(id)
);

