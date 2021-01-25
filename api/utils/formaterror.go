package utils

import (
	"errors"
	"strings"
)

func FormatError(err error) error {
	errString := err.Error()
	if strings.Contains(errString, "nickname") {
		return errors.New("Nickname Already Taken")
	}
	if strings.Contains(errString, "email") {
		return errors.New("Email Already Taken")
	}
	if strings.Contains(errString, "title") {
		return errors.New("Title Already Taken")
	}
	if strings.Contains(errString, "hashedPassword") {
		return errors.New("Incorrect Password")
	}
	return errors.New("Incorrect Details")
}
