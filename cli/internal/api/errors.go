package api

import "errors"

var ErrUnauthenticated = errors.New("unauthenticated")
var ErrPreconditionFailed = errors.New("precondition failed")
