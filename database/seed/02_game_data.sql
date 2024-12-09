-- Seed data for game-related tables

-- Insert test accounts
INSERT INTO accounts (username) 
VALUES 
    ('player1'),
    ('player2'),
    ('player3'),
    ('player4'),
    ('player5')
ON CONFLICT (username) DO NOTHING;

-- Insert characters for each account
-- Class IDs: 0-8 (as per constraint)
INSERT INTO characters (acc_id, class_id)
SELECT 
    acc_id,
    FLOOR(RANDOM() * 8)::int as class_id
FROM accounts
WHERE NOT EXISTS (
    SELECT 1 FROM characters WHERE characters.acc_id = accounts.acc_id
);

-- Insert initial scores for each character
INSERT INTO scores (char_id, reward_score)
SELECT 
    char_id,
    FLOOR(RANDOM() * 1000)::int as reward_score
FROM characters
WHERE NOT EXISTS (
    SELECT 1 FROM scores WHERE scores.char_id = characters.char_id
);
