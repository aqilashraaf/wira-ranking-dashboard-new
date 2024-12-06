-- Count total records in each table
SELECT 
    'accounts' as table_name, COUNT(*) as count 
FROM accounts
UNION ALL
SELECT 
    'characters' as table_name, COUNT(*) as count 
FROM characters
UNION ALL
SELECT 
    'scores' as table_name, COUNT(*) as count 
FROM scores;

-- Sample of accounts with their characters and scores
SELECT 
    a.username,
    c.class_id,
    COUNT(s.score_id) as num_scores,
    MIN(s.reward_score) as min_score,
    MAX(s.reward_score) as max_score,
    AVG(s.reward_score) as avg_score
FROM accounts a
JOIN characters c ON a.acc_id = c.acc_id
JOIN scores s ON c.char_id = s.char_id
GROUP BY a.username, c.class_id
LIMIT 10;
