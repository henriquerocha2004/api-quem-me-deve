package zipcode

import (
	"testing"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		name     string
		cep      string
		expected error
	}{
		{
			name:     "Valid CEP",
			cep:      "12345678",
			expected: nil,
		},
		{
			name:     "Valid CEP with formatting",
			cep:      "12345-678",
			expected: nil,
		},
		{
			name:     "Invalid CEP - Starts with zero",
			cep:      "01234567",
			expected: ErrInvalidCEP,
		},
		{
			name:     "Invalid CEP - Too short",
			cep:      "1234567",
			expected: ErrInvalidLength,
		},
		{
			name:     "Invalid CEP - Too long",
			cep:      "123456789",
			expected: ErrInvalidLength,
		},
		{
			name:     "Invalid CEP - Non-numeric",
			cep:      "1234567a",
			expected: ErrInvalidLength,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.cep)
			if err != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	testCases := []struct {
		name        string
		cep         string
		expectedCEP string
		expectedErr error
	}{
		{
			name:        "Format Valid CEP",
			cep:         "12345678",
			expectedCEP: "12345-678",
			expectedErr: nil,
		},
		{
			name:        "Format CEP with Formatting",
			cep:         "12345-678",
			expectedCEP: "12345-678",
			expectedErr: nil,
		},
		{
			name:        "Format Invalid CEP",
			cep:         "01234567",
			expectedCEP: "",
			expectedErr: ErrInvalidCEP,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formattedCEP, err := Format(tc.cep)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			if formattedCEP != tc.expectedCEP {
				t.Errorf("Expected CEP %s, got %s", tc.expectedCEP, formattedCEP)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	testCases := []struct {
		name        string
		cep         string
		expectedCEP string
	}{
		{
			name:        "Normalize CEP with Formatting",
			cep:         "12345-678",
			expectedCEP: "12345678",
		},
		{
			name:        "Normalize CEP without Formatting",
			cep:         "12345678",
			expectedCEP: "12345678",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalizedCEP := Normalize(tc.cep)
			if normalizedCEP != tc.expectedCEP {
				t.Errorf("Expected normalized CEP %s, got %s", tc.expectedCEP, normalizedCEP)
			}
		})
	}
}

func BenchmarkValidate(b *testing.B) {
	cep := "12345678"
	for i := 0; i < b.N; i++ {
		Validate(cep)
	}
}
