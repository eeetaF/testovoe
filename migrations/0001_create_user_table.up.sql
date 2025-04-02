CREATE TABLE IF NOT EXISTS usr (
                                      id              BIGSERIAL PRIMARY KEY,
                                      name            VARCHAR(30) UNIQUE NOT NULL,
                                      password        VARCHAR(255) NOT NULL,
                                      created_unix    BIGINT NOT NULL,
                                      updated_unix    BIGINT NOT NULL,
                                      referal_code    VARCHAR(7) UNIQUE NOT NULL,
                                      twitter_points  BIGINT DEFAULT 0 NOT NULL,
                                      telegram_points BIGINT DEFAULT 0 NOT NULL,
                                      sum_points      BIGINT GENERATED ALWAYS AS (twitter_points + telegram_points) STORED
);

CREATE INDEX IF NOT EXISTS idx_user_id           ON usr (id);
CREATE INDEX IF NOT EXISTS idx_user_name         ON usr (name);
CREATE INDEX IF NOT EXISTS idx_user_ref_code     ON usr (referal_code);
CREATE INDEX IF NOT EXISTS idx_user_sum_points   ON usr (sum_points);
