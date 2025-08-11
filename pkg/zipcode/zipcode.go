package zipcode

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrInvalidCEP    = errors.New("invalid CEP")
	ErrInvalidFormat = errors.New("invalid CEP format")
	ErrInvalidLength = errors.New("CEP must be 8 digits long")
)

func Validate(cep string) error {
	cep = regexp.MustCompile(`\D`).ReplaceAllString(cep, "")
	if len(cep) != 8 {
		return ErrInvalidLength
	}

	if !regexp.MustCompile(`^\d+$`).MatchString(cep) {
		return ErrInvalidFormat
	}

	primeiroDigito := cep[0]
	if primeiroDigito == '0' {
		return ErrInvalidCEP
	}

	return nil
}

func Format(cep string) (string, error) {
	cep = regexp.MustCompile(`\D`).ReplaceAllString(cep, "")
	if err := Validate(cep); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", cep[:5], cep[5:]), nil
}

func Normalize(cep string) string {
	return regexp.MustCompile(`\D`).ReplaceAllString(cep, "")
}
