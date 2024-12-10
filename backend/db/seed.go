package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultNumUsers = 5000 // This will generate 100,000+ total records due to multiple characters per user
)

var (
	malayPrefixes = []string{
		"Sang", "Si", "Tun", "Tok", "Megat", "Nik", "Wan", "Raja", "Putera", "Awang",
	}

	malayNames = []string{
		"Tuah", "Jebat", "Lekir", "Lekiu", "Kasturi", "Setia", "Perkasa", "Pahlawan", "Laksamana", "Hulubalang",
		"Satria", "Wira", "Kesuma", "Sakti", "Gagah", "Berani", "Laksana", "Andika", "Mahkota", "Bijaksana",
	}

	malayTitles = []string{
		"Pendekar", "Hulubalang", "Laksamana", "Panglima", "Satria", "Wira", "Kesatria", "Pahlawan", "Perwira", "Jaguh",
	}

	characterClasses = []int{1, 2, 3, 4, 5, 6, 7, 8}
)

func generateUniqueName(id int) string {
	prefix := malayPrefixes[rand.Intn(len(malayPrefixes))]
	name := malayNames[rand.Intn(len(malayNames))]
	title := malayTitles[rand.Intn(len(malayTitles))]
	return fmt.Sprintf("%s %s %s %d", prefix, name, title, id)
}

func generateEmail(username string) string {
	// Replace spaces with dots and make lowercase
	email := strings.ToLower(strings.ReplaceAll(username, " ", "."))
	return fmt.Sprintf("%s@wira-ranking.com", email)
}

// SeedData seeds initial data into the database
func SeedData(db *sql.DB) error {
	rand.Seed(time.Now().UnixNano())

	// Get number of users to generate from environment variable
	numUsers := os.Getenv("SEED_NUM_USERS")
	if numUsers == "" {
		numUsers = strconv.Itoa(defaultNumUsers)
	}
	n, err := strconv.Atoi(numUsers)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare statements
	stmtAccount, err := tx.Prepare(`
		INSERT INTO accounts (username, email) 
		VALUES ($1, $2) 
		RETURNING acc_id`)
	if err != nil {
		return err
	}

	stmtCharacter, err := tx.Prepare(`
		INSERT INTO characters (acc_id, class_id) 
		VALUES ($1, $2) 
		RETURNING char_id`)
	if err != nil {
		return err
	}

	stmtScore, err := tx.Prepare(`
		INSERT INTO scores (char_id, reward_score) 
		VALUES ($1, $2)`)
	if err != nil {
		return err
	}

	// Generate data
	for i := 0; i < n; i++ {
		// Generate account
		username := generateUniqueName(i + 1)  // Add 1 to avoid starting from 0
		email := generateEmail(username)
		
		var accID int
		err = stmtAccount.QueryRow(username, email).Scan(&accID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Generate 3-8 characters per account
		numCharacters := rand.Intn(6) + 3
		for j := 0; j < numCharacters; j++ {
			// Generate character with class_id between 1-8
			classID := characterClasses[rand.Intn(len(characterClasses))]
			var charID int
			err = stmtCharacter.QueryRow(accID, classID).Scan(&charID)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Generate 5-10 scores per character
			numScores := rand.Intn(6) + 5
			for k := 0; k < numScores; k++ {
				// Generate realistic score based on normal distribution
				score := int(rand.NormFloat64()*1000 + 5000)
				if score < 0 {
					score = 0
				}

				_, err = stmtScore.Exec(charID, score)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		if i%100 == 0 {
			log.Printf("Generated data for %d users...\n", i)
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Printf("Data generation completed successfully!\n")
	log.Printf("Generated:\n- %d users\n- ~%d characters\n- ~%d scores\n",
		n,
		n*5, // average 5 characters per user
		n*5*7) // average 7 scores per character
	return nil
}
