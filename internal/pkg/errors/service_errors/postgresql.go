package service_errors

import "github.com/warehouse/mail-service/internal/pkg/errors"

var (
	DatabaseError    = &errors.Error{Code: 500, Reason: "database failed"}
	DatabaseErrorRaw = errors.New("database failed")
)
