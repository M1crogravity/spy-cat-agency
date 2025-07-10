CREATE TABLE IF NOT EXISTS tokens (
  hash bytea PRIMARY KEY,
  user_id bigint NOT NULL,
  user_type text NOT NULL,
  expiry timestamp(0) with time zone NOT NULL,
  scope text NOT NULL
);
