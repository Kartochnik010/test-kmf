CREATE TABLE IF NOT EXISTS r_currency (
    id SERIAL PRIMARY KEY,
    title VARCHAR(60) NOT NULL,
    code VARCHAR(3) NOT NULL,
    value float NOT NULL,
    a_date VARCHAR(10) NOT NULL
);
