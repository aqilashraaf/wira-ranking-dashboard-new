package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"wira-dashboard/models"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// GetRankings returns the rankings with pagination
func (h *Handler) GetRankings(c *gin.Context) {
	// Log request
	log.Printf("Received rankings request from: %s with query params: %v", c.Request.RemoteAddr, c.Request.URL.Query())

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Printf("Error parsing page parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if err != nil {
		log.Printf("Error parsing per_page parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid per_page parameter"})
		return
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	// Get class_id from query, default to 0 (all classes)
	classIDStr := c.DefaultQuery("class_id", "0")
	log.Printf("Raw class_id parameter: %s", classIDStr)
	
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		log.Printf("Error parsing class_id parameter: %v", err)
		classID = 0 // Default to all classes if invalid
	}

	if classID < 0 || classID > 8 {
		log.Printf("Invalid class_id value: %d, defaulting to 0", classID)
		classID = 0
	}

	// Log query parameters
	log.Printf("Processed query params - page: %d, perPage: %d, classID: %d", page, perPage, classID)

	var rows *sql.Rows
	var total int

	if classID > 0 {
		err = h.db.QueryRow(`
			SELECT COUNT(*) FROM (
				SELECT DISTINCT ON (a.username) a.username
				FROM accounts a
				JOIN characters c ON a.acc_id = c.acc_id
				JOIN scores s ON c.char_id = s.char_id
				WHERE c.class_id = $1
			) AS unique_users`, classID).Scan(&total)

		if err != nil {
			log.Printf("Error getting total count: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		rows, err = h.db.Query(`
			WITH RankedScores AS (
				SELECT 
					a.username,
					c.class_id,
					MAX(s.reward_score) as highest_score,
					DENSE_RANK() OVER (ORDER BY MAX(s.reward_score) DESC) as rank
				FROM accounts a
				JOIN characters c ON a.acc_id = c.acc_id
				JOIN scores s ON c.char_id = s.char_id
				WHERE c.class_id = $1
				GROUP BY a.username, c.class_id
			)
			SELECT username, class_id, highest_score, rank
			FROM RankedScores
			ORDER BY rank
			LIMIT $2 OFFSET $3`, classID, perPage, offset)
	} else {
		err = h.db.QueryRow(`
			SELECT COUNT(*) FROM (
				SELECT DISTINCT ON (a.username) a.username
				FROM accounts a
				JOIN characters c ON a.acc_id = c.acc_id
				JOIN scores s ON c.char_id = s.char_id
			) AS unique_users`).Scan(&total)

		if err != nil {
			log.Printf("Error getting total count: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		rows, err = h.db.Query(`
			WITH RankedScores AS (
				SELECT 
					a.username,
					c.class_id,
					MAX(s.reward_score) as highest_score,
					DENSE_RANK() OVER (ORDER BY MAX(s.reward_score) DESC) as rank
				FROM accounts a
				JOIN characters c ON a.acc_id = c.acc_id
				JOIN scores s ON c.char_id = s.char_id
				GROUP BY a.username, c.class_id
			)
			SELECT username, class_id, highest_score, rank
			FROM RankedScores
			ORDER BY highest_score DESC
			LIMIT $1 OFFSET $2`, perPage, offset)
	}

	if err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var rankings []models.RankingResponse
	for rows.Next() {
		var r models.RankingResponse
		if err := rows.Scan(&r.Username, &r.ClassID, &r.HighestScore, &r.Rank); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		rankings = append(rankings, r)
	}

	// Log response
	log.Printf("Sending response with %d rankings", len(rankings))

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    rankings,
	})
}

// SearchRankings searches for players by username with pagination
func (h *Handler) SearchRankings(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Printf("Error parsing page parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if err != nil {
		log.Printf("Error parsing per_page parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid per_page parameter"})
		return
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	// Get class_id from query, default to 0 (all classes)
	classIDStr := c.DefaultQuery("class_id", "0")
	log.Printf("Raw class_id parameter: %s", classIDStr)
	
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		log.Printf("Error parsing class_id parameter: %v", err)
		classID = 0 // Default to all classes if invalid
	}

	if classID < 0 || classID > 8 {
		log.Printf("Invalid class_id value: %d, defaulting to 0", classID)
		classID = 0
	}

	// Log query parameters
	log.Printf("Processed query params - page: %d, perPage: %d, classID: %d, username: %s", page, perPage, classID, username)

	query := `
		WITH RankedScores AS (
			SELECT 
				a.username,
				c.class_id,
				MAX(s.reward_score) as highest_score,
				DENSE_RANK() OVER (ORDER BY MAX(s.reward_score) DESC) as rank
			FROM accounts a
			JOIN characters c ON a.acc_id = c.acc_id
			JOIN scores s ON c.char_id = s.char_id
			WHERE LOWER(a.username) LIKE LOWER($1)
			GROUP BY a.username, c.class_id
		)
		SELECT username, class_id, highest_score, rank
		FROM RankedScores
	`

	var rows *sql.Rows
	var total int

	if classID > 0 {
		err = h.db.QueryRow(`
			SELECT COUNT(*) FROM (
				SELECT DISTINCT ON (a.username) a.username
				FROM accounts a
				JOIN characters c ON a.acc_id = c.acc_id
				JOIN scores s ON c.char_id = s.char_id
				WHERE LOWER(a.username) LIKE LOWER($1) AND c.class_id = $2
			) AS unique_users`, "%"+username+"%", classID).Scan(&total)
		
		if err != nil {
			log.Printf("Error getting total count: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		rows, err = h.db.Query(query + ` WHERE class_id = $2 ORDER BY rank LIMIT $3 OFFSET $4`, "%"+username+"%", classID, perPage, offset)
	} else {
		err = h.db.QueryRow(`
			SELECT COUNT(*) FROM (
				SELECT DISTINCT ON (a.username) a.username
				FROM accounts a
				JOIN characters c ON a.acc_id = c.acc_id
				JOIN scores s ON c.char_id = s.char_id
				WHERE LOWER(a.username) LIKE LOWER($1)
			) AS unique_users`, "%"+username+"%").Scan(&total)
		
		if err != nil {
			log.Printf("Error getting total count: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		rows, err = h.db.Query(query + ` ORDER BY rank LIMIT $2 OFFSET $3`, "%"+username+"%", perPage, offset)
	}

	if err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var rankings []models.RankingResponse
	for rows.Next() {
		var r models.RankingResponse
		if err := rows.Scan(&r.Username, &r.ClassID, &r.HighestScore, &r.Rank); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		rankings = append(rankings, r)
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    rankings,
	})
}
