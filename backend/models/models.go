package models

type Account struct {
	AccID    int    `json:"acc_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Character struct {
	CharID  int `json:"char_id"`
	AccID   int `json:"acc_id"`
	ClassID int `json:"class_id"`
}

type Score struct {
	ScoreID      int `json:"score_id"`
	CharID       int `json:"char_id"`
	RewardScore  int `json:"reward_score"`
}

type RankingResponse struct {
	Username     string `json:"username"`
	ClassID      int    `json:"class_id"`
	HighestScore int    `json:"highest_score"`
	Rank         int    `json:"rank"`
}

type PaginatedResponse struct {
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
	Data    interface{} `json:"data"`
}
