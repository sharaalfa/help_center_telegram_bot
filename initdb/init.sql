CREATE TABLE tickets (
                         id SERIAL PRIMARY KEY,
                         department VARCHAR(255) NOT NULL,
                         title VARCHAR(255) NOT NULL,
                         description TEXT NOT NULL,
                         client_id BIGINT NOT NULL
);