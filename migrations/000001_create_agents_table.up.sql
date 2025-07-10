CREATE TABLE IF NOT EXISTS agents(
  id bigserial PRIMARY KEY,
  name text UNIQUE NOT NULL,
  password_hash bytea NOT NULL
)
