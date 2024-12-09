package db

import (
	"database/sql"
	"log"
)

// SeedData seeds initial data into the database
func SeedData(db *sql.DB) error {
	// Seed test accounts
	accounts := []string{"player1", "player2", "player3", "player4", "player5"}
	for _, username := range accounts {
		_, err := db.Exec(`
			INSERT INTO accounts (username) 
			VALUES ($1)
			ON CONFLICT (username) DO NOTHING
		`, username)
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
