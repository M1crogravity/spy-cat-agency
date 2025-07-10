CREATE TABLE IF NOT EXISTS targets (
  id bigint,
  name text NOT NULL,
  country text NOT NULL,
  notes text NOT NULL,
  state text NOT NULL,
  mission_id bigint NULL REFERENCES missions ON DELETE CASCADE,
  PRIMARY KEY (mission_id, id)
);
