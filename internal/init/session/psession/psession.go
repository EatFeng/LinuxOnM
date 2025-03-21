package psession

import (
	"LinuxOnM/internal/init/cache/badger_db"
	"encoding/json"
	"time"
)

type PSession struct {
	ExpireTime int64 `json:"expire_time"`
	store      *badger_db.Cache
}

func NewPSession(db *badger_db.Cache) *PSession {
	return &PSession{
		store: db,
	}
}

type SessionUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (p *PSession) Get(sessionID string) (SessionUser, error) {
	var result SessionUser
	item, err := p.store.Get(sessionID)
	if err != nil {
		return result, err
	}
	_ = json.Unmarshal(item, &result)
	return result, nil
}

func (p *PSession) Set(sessionID string, user SessionUser, ttlSeconds int) error {
	p.ExpireTime = time.Now().Unix() + int64(ttlSeconds)
	return p.store.SetWithTTL(sessionID, user, time.Second*time.Duration(ttlSeconds))
}

func (p *PSession) Delete(sessionID string) error {
	return p.store.Del(sessionID)
}

func (p *PSession) Clean() error {
	return p.store.Clean()
}
