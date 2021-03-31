package model

import "strings"

type QueryError struct {
	Query string
	Err   error
}

func (qe *QueryError) Error() string {
	return qe.Query + ": " + qe.Err.Error()
}

func (qe *QueryError) Unwrap() error {
	return qe.Err
}

type ValidationError struct {
	errors []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.errors, ",")
}
