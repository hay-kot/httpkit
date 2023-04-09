package main

import (
	"errors"
	"fmt"

	"github.com/hay-kot/safeserve/errtrace"
)

func CreateUser() error {
	err := errors.New("user with id 1 already exists")
	err = fmt.Errorf("wrap: %w", err)
	return errtrace.Wrapf(err, "error writing to database")
}

func RepoNewUser() error {
	err := CreateUser()
	return fmt.Errorf("user repo: %w", err)
}

func ServiceNewUser() error {
	err := RepoNewUser()
	return errtrace.Wrapf(err, "error creating user in database")
}

func printErr(err error) {
	trace := errtrace.TraceString(err)

	fmt.Print(trace)
}

func main() {
	println("\n------- Multiple Traceable errors -------\n")
	err := ServiceNewUser()
	err = errtrace.Wrapf(err, "failed to do something")
	printErr(err)

	println("\n------- Multiple Traceable errors wrapped by generic error -------\n")
	err = fmt.Errorf("failed to do something: %w", ServiceNewUser())
	printErr(err)

	println("\n------- Non-Traceable error -------\n")
	err = errors.New("test")
	err = fmt.Errorf("outer: %w", err)
	printErr(err)

	println("\n------- Multiple Traceable errors, but the first one is not Traceable -------\n")
	err = fmt.Errorf("outer: %w", errtrace.Wrapf(errors.New("inner error"), "failed to do something"))
	printErr(err)
}
