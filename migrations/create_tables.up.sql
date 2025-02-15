-- Creating Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Creating Books table
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author_id INT NOT NULL,  -- Assuming you have an authors table for author_id to reference
    location VARCHAR(100) NOT NULL
);

-- Creating BookIssues table
CREATE TABLE IF NOT EXISTS book_issues (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    book_id INT NOT NULL,
    issue_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP,
    due_date TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- Creating Fines table
CREATE TABLE IF NOT EXISTS fines (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    fine_amount DECIMAL(10, 2) NOT NULL,
    date TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
