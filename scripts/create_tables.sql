CREATE TABLE users (
    id serial PRIMARY KEY,
    name text
);

CREATE TABLE segments (
    id serial PRIMARY KEY,
    slug text UNIQUE,
    description text
);

CREATE TABLE users_segments (
    user_id serial,
    slug text,
    deadline_date timestamp with time zone,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_segment FOREIGN KEY (slug) REFERENCES segments (slug)
);

CREATE TABLE history (
    user_id serial NOT NULL,
    slug text NOT NULL,
    action_date timestamp with time zone NOT NULL,
    operation_type text NOT NULL CHECK (operation_type IN ('ADDING', ' REMOVING')),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_segment FOREIGN KEY (slug) REFERENCES segments (slug)
);
