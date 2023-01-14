package mysql

import "errors"

var (
	RecordNotFoundErr error = errors.New("record not found")
)
