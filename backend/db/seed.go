package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	defaultNumUsers = 5000 // This will generate 100,000+ total records due to multiple characters per user
)

// Malay warrior-themed name prefixes and suffixes
var (
	malayPrefixes = []string{"Hang", "Laksamana", "Tun", "Datuk", "Panglima", "Raja", "Sultan"}
	malayNames    = []string{"Tuah", "Jebat", "Lekir", "Kasturi", "Lekiu", "Pahang", "Melaka", "Perang"}
	malayTitles   = []string{"Perkasa", "Wira", "Pahlawan", "Sakti", "Gagah", "Berani"}
)

func generateMalayWarriorName() string {
	prefix := malayPrefixes[rand.Intn(len(malayPrefixes))]
	name := malayNames[rand.Intn(len(malayNames))]
	title := malayTitles[rand.Intn(len(malayTitles))]
	
	// Add timestamp and random number to ensure uniqueness
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	uniqueNum := rand.Intn(9999)
	return fmt.Sprintf("%s_%s_%s_%d_%d", prefix, name, title, timestamp, uniqueNum)
}

func generateEmail(username string) string {
	return fmt.Sprintf("%s@wira-ranking.com", username)
}

// SeedData seeds initial data into the database
func SeedData(db *sql.DB) error {
	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Get number of users to generate from environment variable
	numUsers := defaultNumUsers
	if envNumUsers := os.Getenv("SEED_NUM_USERS"); envNumUsers != "" {
		if n, err := strconv.Atoi(envNumUsers); err == nil {
			numUsers = n
		}
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
	for i := 0; i < numUsers; i++ {
		// Generate account
		username := generateMalayWarriorName()
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
			// Generate character with class_id between 0-8
			classID := rand.Intn(9)
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
		numUsers,
		numUsers*5, // average 5 characters per user
		numUsers*5*7) // average 7 scores per character
	return nil
}
