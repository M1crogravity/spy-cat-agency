CREATE TABLE IF NOT EXISTS missions (
  id bigserial PRIMARY KEY,
  state text NOT NULL,
  spy_cat_id bigint NULL REFERENCES spy_cats ON DELETE CASCADE
);
