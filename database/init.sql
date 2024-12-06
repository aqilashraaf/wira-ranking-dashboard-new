-- Drop existing tables if they exist
DROP TABLE IF EXISTS scores;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS accounts;

-- Create accounts table
CREATE TABLE accounts (
    acc_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);

-- Create characters table
CREATE TABLE characters (
    char_id SERIAL PRIMARY KEY,
    acc_id INTEGER REFERENCES accounts(acc_id),
    class_id INTEGER CHECK (class_id BETWEEN 1 AND 8) NOT NULL
);

-- Create scores table
CREATE TABLE scores (
    score_id SERIAL PRIMARY KEY,
    char_id INTEGER REFERENCES characters(char_id),
    reward_score INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_characters_acc_id ON characters(acc_id);
CREATE INDEX idx_scores_char_id ON scores(char_id);
CREATE INDEX idx_scores_reward ON scores(reward_score DESC);
