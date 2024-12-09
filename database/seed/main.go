package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"strings"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

const (
	defaultNumUsers = 5000 // This will generate 100,000+ total records due to multiple characters per user
)

// Malay warrior-themed name prefixes and suffixes
var (
	malayPrefixes = []string{"Hang", "Laksamana", "Tun", "Datuk", "Panglima", "Raja", "Sultan"}
	malayNames    = []string{"Tuah", "Jebat", "Lekir", "Kasturi", "Lekiu", "Pahang", "Melaka", "Perang"}
	malayTitles   = []string{"Perkasa", "Wira", "Pahlawan", "Sakti", "Gagah", "Berani"}
	emailDomains  = []string{"gmail.com", "yahoo.com", "hotmail.com", "outlook.com"}
)

func generateMalayWarriorName() (string, string) {
	prefix := malayPrefixes[rand.Intn(len(malayPrefixes))]
	name := malayNames[rand.Intn(len(malayNames))]
	title := malayTitles[rand.Intn(len(malayTitles))]
	
	// Add timestamp and random number to ensure uniqueness
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	uniqueNum := rand.Intn(9999)
	username := fmt.Sprintf("%s_%s_%s_%d_%d", prefix, name, title, timestamp, uniqueNum)
	
	// Generate email using the username
	email := fmt.Sprintf("%s@wira.com", strings.ToLower(username))
	
	return username, email
}

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Get database configuration from environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5434"
	}
	
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "aqash18"
	}
	
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "wira_dashboard"
	}

	// Get number of users to generate from environment variable
	numUsers := defaultNumUsers
	if envNumUsers := os.Getenv("SEED_NUM_USERS"); envNumUsers != "" {
		if n, err := strconv.Atoi(envNumUsers); err == nil {
			numUsers = n
		}
	}

	// Connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database. Starting data generation...")

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare statements
	stmtAccount, err := tx.Prepare(`
		INSERT INTO accounts (username, email)
		VALUES ($1, $2)
		RETURNING acc_id`)
	if err != nil {
		log.Fatal(err)
	}

	stmtCharacter, err := tx.Prepare(`
		INSERT INTO characters (acc_id, class_id)
		VALUES ($1, $2)
		RETURNING char_id`)
	if err != nil {
		log.Fatal(err)
	}

	stmtScore, err := tx.Prepare(`
		INSERT INTO scores (char_id, reward_score)
		VALUES ($1, $2)`)
	if err != nil {
		log.Fatal(err)
	}

	// Generate data
	for i := 0; i < numUsers; i++ {
		// Generate account
		username, email := generateMalayWarriorName()
		
		var accID int
		err = stmtAccount.QueryRow(username, email).Scan(&accID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		// Generate 3-8 characters per account
		numCharacters := rand.Intn(6) + 3
		for j := 0; j < numCharacters; j++ {
			// Generate character with class_id between 1-8 (not 0)
			classID := rand.Intn(8) + 1
			var charID int
			err = stmtCharacter.QueryRow(accID, classID).Scan(&charID)
			if err != nil {
				tx.Rollback()
				log.Fatal(err)
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
					log.Fatal(err)
				}
			}
		}

		if i%100 == 0 {
			fmt.Printf("Generated data for %d users...\n", i)
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data generation completed successfully!")
	fmt.Printf("Generated:\n- %d users\n- ~%d characters\n- ~%d scores\n",
		numUsers,
		numUsers*5, // average 5 characters per user
		numUsers*5*7) // average 7 scores per character
}
