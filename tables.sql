-- User table to track users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Book issues table to track books issued to users
CREATE TABLE book_issues (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    book_id INT REFERENCES books(id),
    issue_date DATE NOT NULL,
    return_date DATE,
    due_date DATE NOT NULL
);

-- Fines table to track fines for users
CREATE TABLE fines (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    fine_amount DECIMAL(10, 2) NOT NULL,
    date DATE NOT NULL
);
