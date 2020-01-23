
CREATE TABLE if not exists Users (
   id SERIAL PRIMARY KEY,
   subject        TEXT    NOT NULL UNIQUE,
   email          TEXT    NOT NULL UNIQUE,
   provider          TEXT    NOT NULL UNIQUE,
   last_login  timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE if not exists Tournaments (
   id SERIAL PRIMARY KEY,
   tournament_name TEXT,
   user_id integer NOT NULL,
   tournament_data json,
   creation_date  timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE if not exists Games (
   id SERIAL PRIMARY KEY NOT NULL,
   gameJson json NOT NULL

);


