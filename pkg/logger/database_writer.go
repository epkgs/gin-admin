package logger

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

var TableName = "sys_log"

type LoggerModel struct {
	ID     string    `json:"id" gorm:"size:20;primaryKey;"`            // Trace ID (Unique)
	Time   time.Time `json:"time" gorm:"index;"`                       // Log time
	Level  string    `json:"level" gorm:"size:20;index;"`              // Log level
	Msg    string    `json:"msg" gorm:"size:1024;"`                    // Log message
	Source Source    `json:"source" gorm:"type:text;serializer:json;"` // Source of error

	Attrs map[string]any `json:"meta" gorm:"type:text;serializer:json;"` // Log data
}

func (LoggerModel) TableName() string {
	return TableName
}

func newDBWriter(db *gorm.DB) io.Writer {

	db.AutoMigrate(&LoggerModel{})

	return &DBWriter{
		db: db,
	}
}

type DBWriter struct {
	db *gorm.DB
}

func (a *DBWriter) Write(b []byte) (n int, err error) {

	data := make(map[string]any)
	err = jsoniter.Unmarshal(b, &data)
	if err != nil {
		return
	}

	model := LoggerModel{}

	if v, ok := data["id"]; ok { // trace id
		model.ID = v.(string)
		delete(data, "id")
	} else {
		model.ID = xid.New().String()
	}

	if v, ok := data["time"]; ok { // slog time key
		// 将 2024-01-01T12:00:00.000+08:00 格式化为 time.Time
		model.Time, err = time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return
		}
		delete(data, "time")
	}

	if v, ok := data["level"]; ok { // slog level key
		model.Level = v.(string)
		delete(data, "level")
	}

	if v, ok := data["msg"]; ok { // slog message key
		model.Msg = v.(string)
		delete(data, "msg")
	}

	if v, ok := data["source"]; ok { // slog source key
		if src, ok := v.(map[string]any); ok {
			if file, ok := src["file"].(string); ok {
				model.Source.File = file
			}
			if function, ok := src["function"].(string); ok {
				model.Source.Function = function
			}
			if line, ok := src["line"].(int); ok {
				model.Source.Line = line
			}
		}
		delete(data, "source")
	}

	model.Attrs = data

	err = a.db.Create(&model).Error
	n = len(b)
	return
}
