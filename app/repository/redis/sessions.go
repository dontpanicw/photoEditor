package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

type SessionCashStorage struct {
	rdb *redis.Client
}

func (rs *SessionCashStorage) GetSessionCashStorage() *redis.Client {
	return rs.rdb
}

func NewSessionCashStorage(connStr string, password string, dbNumber int) (*SessionCashStorage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     connStr,  // redis:6379
		Password: password, // no password set
		DB:       dbNumber, // use default DB
	})
	return &SessionCashStorage{rdb}, nil
}

var ctx = context.Background()

func (rs *SessionCashStorage) CreateNewSession(userId int64) (int64, error) {
	rdb := rs.GetSessionCashStorage()
	sessionId, err := rdb.Incr(ctx, "session:id:counter").Result()
	if err != nil {
		log.Println("Ошибка при получении ID:", err)
		return 0, err
	}

	key := strconv.FormatInt(sessionId, 10)

	err = rdb.Set(ctx, key, userId, 0).Err()
	if err != nil {
		log.Println("Ошибка при создании сессии:", err)
		return 0, err
	}
	return sessionId, nil
}

func (rs *SessionCashStorage) GetSession(sessionId int64) error {
	rdb := rs.GetSessionCashStorage()
	key := strconv.FormatInt(sessionId, 10)
	_, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return errors.New("session not found")
	}
	return nil
}

func (rs *SessionCashStorage) GetAllSessions() map[string]string {
	rdb := rs.GetSessionCashStorage()
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		return nil
	}

	sessions := make(map[string]string)

	for _, key := range keys {
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		sessions[key] = val
	}

	return sessions

}
