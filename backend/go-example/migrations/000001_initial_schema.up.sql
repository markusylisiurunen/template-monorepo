BEGIN;

CREATE TABLE IF NOT EXISTS messages (
  message_id BIGSERIAL PRIMARY KEY,
  message_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  message_text TEXT NOT NULL
);

COMMIT;
