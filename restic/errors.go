package restic

import "errors"

var (
	ErrConnectionTimeout = errors.New("restic: respository connection timeout")
)
