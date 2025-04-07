package ram_storage

//import (
//	"errors"
//	"homework-dontpanicw/app/domain"
//	"homework-dontpanicw/app/repository"
//	"sync/atomic"
//)
//
//type UserRepository struct {
//	user      map[int64]*domain.User
//	session   map[int64]*domain.Session
//	userId    atomic.Int64
//	sessionId atomic.Int64
//}
//
//func NewUser() repository.User {
//	return &UserRepository{
//		user:    make(map[int64]*domain.User),
//		session: make(map[int64]*domain.Session),
//	}
//}
//
//func (rs *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
//	for _, user := range rs.user {
//		if user.Username == username {
//			return user, nil
//		}
//	}
//	return nil, errors.New("user not found")
//}
//
//func (rs *UserRepository) RegisterNewUser(username string, password string) error {
//	user, _ := rs.GetUserByUsername(username)
//	if user != nil {
//		return errors.New("user already exists")
//	}
//	id := rs.userId.Add(1)
//	rs.user[id] = &domain.User{
//		Id:       id,
//		Username: username,
//		Password: password,
//	}
//	return nil
//}
//
////func (rs *UserRepository) CreateNewSession(id int64) (int64, error) {
////	sessionId := rs.sessionId.Add(1)
////	rs.session[sessionId] = &domain.Session{
////		SessionId: sessionId,
////		UserId:    id,
////	}
////	return sessionId, nil
////}
//
//func (rs *UserRepository) GetSession(sessionId int64) error {
//	_, ok := rs.session[sessionId]
//	if !ok {
//		return errors.New("session not found")
//	}
//	return nil
//}
//
//func (rs *UserRepository) GetAllSessions() map[int64]*domain.Session {
//	return rs.session
//}
