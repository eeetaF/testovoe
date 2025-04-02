ALTER TABLE usr
    ADD COLUMN referal_points BIGINT DEFAULT 0 NOT NULL;

DROP INDEX IF EXISTS idx_user_sum_points;

ALTER TABLE usr
    DROP COLUMN sum_points;

ALTER TABLE usr
    ADD COLUMN sum_points BIGINT GENERATED ALWAYS AS (twitter_points + telegram_points + referal_points) STORED;

CREATE INDEX idx_user_sum_points ON usr(sum_points);
