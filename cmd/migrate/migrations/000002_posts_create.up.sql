CREATE TABLE
    IF NOT EXISTS posts (
        id bigserial PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        user_id bigint NOT NULL,
        content VARCHAR(255) NOT NULL,
        created_at timestamp(0)
        with
            time zone NOT NULL DEFAULT NOW ()
    )