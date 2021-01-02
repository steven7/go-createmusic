package models

import "strings"

const (
// ErrNotFound is returned when a resource cannot be found
// in the database.
ErrNotFound modelError = "models: resource not found"

// ErrIDInvalid is returned when an invalid ID is provided
// to a method like Delete.
ErrIDInvalid modelError = "models: ID provided was invalid"

// ErrPasswordIncorrect is returned when an invalid password
// is used when attempting to authenticate a user.
ErrPasswordIncorrect modelError = "models: incorrect password provided"

// ErrPasswordTooShort is returned when a user tries to set
// a password that is less than 8 characters long.
ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"

// ErrPasswordRequired is returned when a create is attempted
// without a user password provided.
ErrPasswordRequired modelError = "models: password is required"

// ErrEmailRequired is returned when an email address is
// not provided when creating a user
ErrEmailRequired modelError = "models: email address is required"

// ErrEmailInvalid is returned when an email address provided
// does not match any of our requirements
ErrEmailInvalid modelError = "models: email address is not valid"

// ErrEmailTaken is returned when an update or create is attempted
// with an email address that is already in use.
ErrEmailTaken modelError = "models: email address is already taken"

// ErrRememberRequired is returned when a create or update
// is attempted without a user remember token hash
ErrRememberRequired modelError = "models: remember token is required"

// ErrRememberTooShort is returned when a remember token is
// not at least 32 bytes
ErrRememberTooShort modelError = "models: remember token must be at least 32 bytes"

ErrTokenInvalid modelError = "models: token provided is not valid"

ErrUserIDRequired modelError = "models: user ID is required"

ErrTitleRequired  modelError = "models: title is required"

ErrServiceRequired modelError = "models: sevice is required"

//
// other errors
//

ErrNoRecord modelError = "models: no matching record found"
// Add a new ErrInvalidCredentials error. We'll use this later if a user
// tries to login with an incorrect email address or password.
ErrInvalidCredentials modelError = "models: invalid credentials"
// Add a new ErrDuplicateEmail error. We'll use this later if a user
// tries to signup with an email address that's already in use.
ErrDuplicateEmail modelError = "models: duplicate email"

)


type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}