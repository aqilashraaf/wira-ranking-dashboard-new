package db

import (
	"database/sql"
	"log"
)

// SeedData seeds initial data into the database
func SeedData(db *sql.DB) error {
	// Seed test accounts
	accounts := []struct {
		username string
		email    string
	}{
		{"player1", "player1@wira-ranking.com"},
		{"player2", "player2@wira-ranking.com"},
		{"player3", "player3@wira-ranking.com"},
		{"player4", "player4@wira-ranking.com"},
		{"player5", "player5@wira-ranking.com"},
	}
	for _, acc := range accounts {
		_, err := db.Exec(`
			INSERT INTO accounts (username, email) 
			VALUES ($1, $2)
			ON CONFLICT (username) DO NOTHING
		`, acc.username, acc.email)
		if err != nil {
			return err
		}
	}
	log.Println("Successfully seeded accounts")

	// Seed characters with random class_ids
	_, err := db.Exec(`
		INSERT INTO characters (acc_id, class_id)
		SELECT 
			acc_id,
			FLOOR(RANDOM() * 8)::int as class_id
		FROM accounts
		WHERE NOT EXISTS (
			SELECT 1 FROM characters WHERE characters.acc_id = accounts.acc_id
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Successfully seeded characters")

	// Seed scores with random reward_scores
	_, err = db.Exec(`
		INSERT INTO scores (char_id, reward_score)
		SELECT 
			char_id,
			FLOOR(RANDOM() * 1000)::int as reward_score
		FROM characters
		WHERE NOT EXISTS (
			SELECT 1 FROM scores WHERE scores.char_id = characters.char_id
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Successfully seeded scores")

	return nil
}
