package database

import "errors"

var (
	ErrNoRecord       = errors.New("no record found")
	ErrAuthentication = errors.New("email or password is not correct")
)
