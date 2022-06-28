CREATE TABLE climate_records
(
  id SERIAL PRIMARY KEY,
  day INTEGER NOT NULL,
  perimeter FLOAT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  climate VARCHAR(128) NOT NULL
);

CREATE INDEX day_idx on climate_records(day)
