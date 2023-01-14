package error

import (
	"errors"
	"gorm.io/gorm"
)

func RecordNotFound(db *gorm.DB) bool {
	return errors.Is(db.Error, gorm.ErrRecordNotFound)
}
