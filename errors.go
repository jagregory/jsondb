package jsondb

import "fmt"

type NotFoundError struct {
	EntityId string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("Entity not found #%s", err.EntityId)
}

func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
