package handler

import "time"

// requests DTO
type RegisterReq struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// responses DTO
type UserResp struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"createdAt"`
}
