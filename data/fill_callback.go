package data

import (
	"slices"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	authnMiddleware "github.com/tx7do/kratos-authn/middleware"
	"gorm.io/gorm"
)

var auditExcludeTables = []string{
	"base_log",
}

func safeSetColumn(db *gorm.DB, fieldName string, value interface{}) {
	if !isFieldZero(db, fieldName) {
		return
	}
	db.Statement.SetColumn(fieldName, value)
}

func isFieldZero(db *gorm.DB, fieldName string) bool {
	statement := db.Statement
	if statement == nil {
		return false
	}
	if field := statement.Schema.LookUpField(fieldName); field == nil {
		return false
	}
	return true
}

func fillCreatedFields(db *gorm.DB) {
	var err error
	table := db.Statement.Table
	if slices.Contains(auditExcludeTables, table) {
		return
	}

	var userId int64
	ctx := db.Statement.Context
	if authClaims, ok := authnMiddleware.FromContext(ctx); !ok {
		log.Warnf("context has no user info, use default user id")
	} else {
		userId, err = authClaims.GetInt64("userId")
		if err != nil {
			log.Errorf("get user id failed, use default user id")
		}
	}

	now := time.Now()
	safeSetColumn(db, "CreatedBy", userId)
	safeSetColumn(db, "UpdatedBy", userId)
	safeSetColumn(db, "CreatedAt", now)
	safeSetColumn(db, "UpdatedAt", now)
}

func fillUpdatedFields(db *gorm.DB) {
	var err error
	table := db.Statement.Table
	if slices.Contains(auditExcludeTables, table) {
		return
	}

	var userId int64
	ctx := db.Statement.Context
	if authClaims, ok := authnMiddleware.FromContext(ctx); !ok {
		log.Warnf("context has no user info, use default user id")
	} else {
		userId, err = authClaims.GetInt64("userId")
		if err != nil {
			log.Errorf("get user id failed, use default user id")
		}
	}

	safeSetColumn(db, "UpdatedBy", userId)
	safeSetColumn(db, "UpdatedAt", time.Now())
}

func registerFillCallback(db *gorm.DB) error {
	if err := db.Callback().Create().Before("gorm:before_create").Register("fill_created_fields", fillCreatedFields); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:before_update").Register("fill_updated_fields", fillUpdatedFields); err != nil {
		return err
	}
	return nil
}
