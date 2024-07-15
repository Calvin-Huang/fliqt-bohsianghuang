package util

import (
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"fliqt/config"
)

var tracer = otel.Tracer(GetFileNameFromCaller())

func beforeQuery(db *gorm.DB) {
	_, span := tracer.Start(db.Statement.Context, "gorm.query")
	db.InstanceSet("otel:span", span)
}

func afterQuery(db *gorm.DB) {
	instance, ok := db.InstanceGet("otel:span")
	if !ok || instance == nil {
		return
	}
	span, ok := instance.(trace.Span)
	if !ok || span == nil {
		return
	}
	endTrace(db, span)
}

func endTrace(db *gorm.DB, span trace.Span) {
	if !span.IsRecording() {
		return
	}

	statusCode := codes.Ok
	message := "OK"
	if db.Error != nil {
		err := db.Error
		if err == gorm.ErrRecordNotFound {
			statusCode = codes.Error
		} else {
			statusCode = codes.Error
		}

		message = err.Error()
	}
	span.SetStatus(statusCode, message)
	// print sql statement and fuzzy value
	span.AddEvent("sql", trace.WithAttributes(attribute.String("log", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...))))
	span.End()
}

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.GetDBDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		db = db.Debug()
	} else {
		db.Logger = logger.Default.LogMode(logger.Silent)
	}

	sqlDB.SetMaxIdleConns(cfg.DBMaxIdle)
	sqlDB.SetMaxOpenConns(cfg.DBMaxConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBMaxLifeTime))

	db.Callback().Create().After("*").Register("sqlCapture:after_query", afterQuery)
	db.Callback().Delete().After("*").Register("sqlCapture:after_query", afterQuery)
	db.Callback().Update().After("*").Register("sqlCapture:after_query", afterQuery)

	db.Callback().Query().Before("gorm:query").Register("sqlCapture:before_query", beforeQuery)
	db.Callback().Query().After("gorm:query").Register("sqlCapture:after_query", afterQuery)

	db.Callback().Row().Before("gorm:row_query").Register("sqlCapture:before_query", beforeQuery)
	db.Callback().Row().After("gorm:row_query").Register("sqlCapture:after_query", afterQuery)

	return db, nil
}
