package repository

type Session interface {
	CreateNewSession(userId int64) (int64, error)
	GetAllSessions() map[string]string
	GetSession(sessionId int64) error
}
