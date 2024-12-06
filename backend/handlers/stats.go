package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type ClassStats struct {
	ClassID      int     `json:"class_id"`
	PlayerCount  int     `json:"player_count"`
	AvgScore     float64 `json:"average_score"`
	HighestScore int     `json:"highest_score"`
	LowestScore  int     `json:"lowest_score"`
}

// GetClassStats returns statistics for each class or a specific class
func (h *Handler) GetClassStats(c *gin.Context) {
	query := `
		SELECT 
			c.class_id,
			COUNT(DISTINCT c.char_id) as player_count,
			ROUND(AVG(s.reward_score)::numeric, 2) as average_score,
			MAX(s.reward_score) as highest_score,
			MIN(s.reward_score) as lowest_score
		FROM characters c
		JOIN scores s ON c.char_id = s.char_id
		GROUP BY c.class_id
		ORDER BY c.class_id`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var stats []ClassStats
	for rows.Next() {
		var stat ClassStats
		if err := rows.Scan(
			&stat.ClassID,
			&stat.PlayerCount,
			&stat.AvgScore,
			&stat.HighestScore,
			&stat.LowestScore,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		stats = append(stats, stat)
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
