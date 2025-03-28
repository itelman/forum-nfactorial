package users

import "regexp"

var (
	usernameRX = regexp.MustCompile(`^[a-zA-Z]{5,}([._]{0,1}[a-zA-Z0-9]{2,})*$`)
	passwordRX = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
)

const (
	nameMinLen = 5
	nameMaxLen = 30

	pwdMinLen = 6
	pwdMaxLen = 20
)
