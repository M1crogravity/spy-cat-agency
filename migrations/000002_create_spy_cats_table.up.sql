CREATE TABLE IF NOT EXISTS spy_cats(
  id bigserial PRIMARY KEY,
  name text UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  years_of_experience integer NOT NULL,
  breed text NOT NULL,
  salary double precision NOT NULL
)
