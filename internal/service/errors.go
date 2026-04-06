package service

import "errors"

// Sentinel errors for fix session lifecycle (use with errors.Is).
var (
	ErrNoFixInProgress  = errors.New("no fix in progress")
	ErrFixAlreadyActive = errors.New("a fix is already in progress")
)
