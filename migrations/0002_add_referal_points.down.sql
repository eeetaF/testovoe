ALTER TABLE usr
    DROP COLUMN referal_points;

DROP INDEX IF EXISTS idx_user_sum_points;

ALTER TABLE usr
    DROP COLUMN sum_points;

ALTER TABLE usr
    ADD COLUMN sum_points BIGINT GENERATED ALWAYS AS (twitter_points + telegram_points) STORED;

CREATE INDEX idx_user_sum_points ON usr(sum_points);
