CREATE TYPE task_status AS ENUM ('TO_DO', 'IN_PROGRESS', 'DONE');

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    status task_status NOT NULL DEFAULT 'TO_DO',
    deleted_at TIMESTAMP NULL
);