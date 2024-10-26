package custom_errors

import (
	"errors"
)

var (
	ErrStudentNotFound = errors.New("student does not exist")
	ErrStudentExists   = errors.New("student already exists")
	ErrGroupNotFound   = errors.New("group does not exist")
	ErrGroupExists     = errors.New("group already exists")
)
