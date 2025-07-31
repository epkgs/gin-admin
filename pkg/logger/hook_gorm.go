package logger

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Logger struct {
	ID        string    `json:"id" gorm:"size:20;primaryKey;"` // Unique ID
	Level     string    `json:"level" gorm:"size:20;index;"`   // Log level
	Message   string    `json:"message" gorm:"size:1024;"`     // Log message
	CreatedAt time.Time `json:"createdAt" gorm:"index;"`       // Create time
	TraceID   string    `json:"traceId" gorm:"size:64;index;"` // Trace ID
	UserID    string    `json:"userId" gorm:"size:20;index;"`  // User ID
	Tag       string    `json:"tag" gorm:"size:32;index;"`     // Log tag
	Stack     string    `json:"stack" gorm:"type:text;"`       // Error stack

	Meta map[string]any `json:"meta" gorm:"type:text;serializer:json;"` // Log data
}

func NewGormHook(db *gorm.DB) *GormHook {
	err := db.AutoMigrate(new(Logger))
	if err != nil {
		panic(err)
	}

	return &GormHook{
		db: db,
	}
}

// Gorm Logger Hook
type GormHook struct {
	db *gorm.DB
}

func (h *GormHook) Exec(extra map[string]string, b []byte) error {
	msg := &Logger{
		ID: xid.New().String(),
	}
	data := make(map[string]any)
	err := jsoniter.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	if v, ok := data["ts"]; ok { // zap time key
		msg.CreatedAt = time.UnixMilli(int64(v.(float64)))
		delete(data, "ts")
	}
	if v, ok := data["msg"]; ok { // zap message key
		msg.Message = v.(string)
		delete(data, "msg")
	}
	if v, ok := data["level"]; ok { // zap level key
		msg.Level = v.(string)
		delete(data, "level")
	}
	if v, ok := data[key_traceID]; ok { // traceId in context
		msg.TraceID = v.(string)
		delete(data, key_traceID)
	}
	if v, ok := data[key_userID]; ok { // userId in context
		msg.UserID = v.(string)
		delete(data, key_userID)
	}
	if v, ok := data[key_tag]; ok { // tag in context
		msg.Tag = v.(string)
		delete(data, key_tag)
	}
	if v, ok := data[key_stack]; ok { // tag in context
		msg.Stack = v.(string)
		delete(data, key_stack)
	}
	delete(data, "caller")

	for k, v := range extra {
		data[k] = v
	}

	msg.Meta = data

	return h.db.Create(msg).Error
}

func (h *GormHook) Close() error {
	db, err := h.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
