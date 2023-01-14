package mysql

import (
	"database/sql"
	"database/sql/driver"
)

type Serializable interface {
	driver.Valuer
	sql.Scanner
}
