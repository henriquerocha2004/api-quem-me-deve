package document

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	ErrInvalidDocument    = errors.New("invalid document")
	ErrDocumentNotNumeric = errors.New("document must contain only numbers")
	ErrRepeatedDigits     = errors.New("document with repeated digits is invalid")
)

type Document string

func (d *Document) Validate() error {
	err := ValidateCPF(string(*d))
	if err == nil {
		return nil
	}

	return ValidateCNPJ(string(*d))
}

func ValidateCPF(cpf string) error {
	cpf = regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return ErrInvalidDocument
	}

	if isRepeatedDigits(cpf) {
		return ErrRepeatedDigits
	}

	sum := 0
	for i := range 9 {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (10 - i)
	}
	firstCheck := 11 - (sum % 11)
	if firstCheck >= 10 {
		firstCheck = 0
	}

	firstCheckDigit, _ := strconv.Atoi(string(cpf[9]))
	if firstCheck != firstCheckDigit {
		return ErrInvalidDocument
	}

	sum = 0
	for i := range 10 {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (11 - i)
	}
	secondCheck := 11 - (sum % 11)
	if secondCheck >= 10 {
		secondCheck = 0
	}

	secondCheckDigit, _ := strconv.Atoi(string(cpf[10]))
	if secondCheck != secondCheckDigit {
		return ErrInvalidDocument
	}

	return nil
}

func ValidateCNPJ(cnpj string) error {
	cnpj = regexp.MustCompile(`\D`).ReplaceAllString(cnpj, "")

	if len(cnpj) != 14 {
		return ErrInvalidDocument
	}

	if isRepeatedDigits(cnpj) {
		return ErrRepeatedDigits
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	sum := 0
	for i := range 12 {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weights1[i]
	}
	firstCheck := sum % 11
	firstCheck = 11 - firstCheck
	if firstCheck >= 10 {
		firstCheck = 0
	}

	firstCheckDigit, _ := strconv.Atoi(string(cnpj[12]))
	if firstCheck != firstCheckDigit {
		return ErrInvalidDocument
	}

	sum = 0
	for i := range 13 {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weights2[i]
	}
	secondCheck := sum % 11
	secondCheck = 11 - secondCheck
	if secondCheck >= 10 {
		secondCheck = 0
	}

	secondCheckDigit, _ := strconv.Atoi(string(cnpj[13]))
	if secondCheck != secondCheckDigit {
		return ErrInvalidDocument
	}

	return nil
}

func isRepeatedDigits(document string) bool {
	firstDigit := document[0]
	for i := 1; i < len(document); i++ {
		if document[i] != firstDigit {
			return false
		}
	}
	return true
}
