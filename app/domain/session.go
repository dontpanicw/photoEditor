package domain

type Session struct {
	UserId    int64 `json:"user_id"`
	SessionId int64 `json:"session_id"`
}
