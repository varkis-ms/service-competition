CREATE TABLE IF NOT EXISTS competition
(
    id                  BIGSERIAL   PRIMARY KEY,
    user_id             INTEGER     NOT NULL,
    title               VARCHAR     NOT NULL,
    description         VARCHAR     NOT NULL,
    dataset_title       VARCHAR     NOT NULL,
    dataset_description VARCHAR     NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (title)
);

CREATE TABLE IF NOT EXISTS leaderboard
(
    id             BIGSERIAL   PRIMARY KEY,
    competition_id INTEGER     NOT NULL REFERENCES competition (id),
    user_id        INTEGER     NOT NULL,
    score          FLOAT       NOT NULL DEFAULT 0,
    run_time       INTERVAL    NOT NULL DEFAULT '0 days',
    added_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    status         VARCHAR     NOT NULL DEFAULT 'in queue'
);
