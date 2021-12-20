CREATE TABLE if not exists Accounts (
   id SERIAL NOT NULL,
   admin_user_id    integer,
   account_name           TEXT    NOT NULL UNIQUE,
   PRIMARY KEY(id)
);


CREATE TABLE if not exists Images (
   id SERIAL NOT NULL,
   format         TEXT    NOT NULL,
   imageData      BYTEA,
   PRIMARY KEY(id)
);



CREATE TABLE if not exists Players (
   id SERIAL NOT NULL,
   account_id integer NOT NULL,
   player_name    TEXT    NOT NULL,
   email          TEXT,
   phone          TEXT,
   paid           integer NOT NULL DEFAULT 0,
   image_id INTEGER NOT NULL DEFAULT 0,
   PRIMARY KEY(id)
);

CREATE TABLE if not exists PlayerResults (
   player_id integer NOT NULL,
   opponent_id integer NOT NULL,
   game_date DATE NOT NULL DEFAULT CURRENT_DATE,
   win integer NOT NULL
);


CREATE TABLE if not exists Users (
   id SERIAL PRIMARY KEY NOT NULL,
   account_id integer NOT NULL,
   subject        TEXT    NOT NULL UNIQUE,
   email          TEXT,
   provider       TEXT,
   given_name     TEXT,
   picture_url    TEXT,
   last_login  timestamp without time zone DEFAULT CURRENT_TIMESTAMP

);

select setval(pg_get_serial_sequence('users', 'id'), 2);


CREATE TABLE if not exists Tournaments (
   id SERIAL PRIMARY KEY,
   account_id integer NOT NULL,
   user_id integer NOT NULL DEFAULT 0 references Users(id),
   tournament_name TEXT,
   subject TEXT NOT NULL references Users(subject),
   tournament_data json,
   start_time timestamp NOT NULL DEFAULT TIMESTAMP 'epoch',
   end_time timestamp NOT NULL DEFAULT TIMESTAMP 'epoch',
   final_game integer NOT NULL DEFAULT 0,
   tournament_state TEXT NOT NULL DEFAULT 'none',
   participant_count integer NOT NULL DEFAULT '0',
   creation_date  timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE if not exists Preferences (
   id integer PRIMARY KEY,
   subject TEXT NOT NULL DEFAULT '' references Users(subject),
   tournament_id integer,
   preferences_data json
);


insert into Accounts (admin_user_id, account_name) values (0, 'master') ;


