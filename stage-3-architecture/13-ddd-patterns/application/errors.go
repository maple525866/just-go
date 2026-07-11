// Package application orchestrates ordering use cases around the domain model.
package application

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrOrderConflict = errors.New("order version conflict")
	ErrEventPublish  = errors.New("event publication failed")
	ErrDependency    = errors.New("application dependency is required")
)
