package document

import (
	"testing"
)

func TestValidateCPF(t *testing.T) {
	testCases := []struct {
		name     string
		cpf      string
		expected error
	}{
		{
			name:     "Valid CPF",
			cpf:      "52998224725",
			expected: nil,
		},
		{
			name:     "Valid CPF with formatting",
			cpf:      "529.982.247-25",
			expected: nil,
		},
		{
			name:     "Invalid CPF - Wrong check digits",
			cpf:      "52998224726",
			expected: ErrInvalidDocument,
		},
		{
			name:     "Invalid CPF - Repeated digits",
			cpf:      "11111111111",
			expected: ErrRepeatedDigits,
		},
		{
			name:     "Invalid CPF - Incorrect length",
			cpf:      "5299822472",
			expected: ErrInvalidDocument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCPF(tc.cpf)
			if err != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestValidateCNPJ(t *testing.T) {
	testCases := []struct {
		name     string
		cnpj     string
		expected error
	}{
		{
			name:     "Valid CNPJ",
			cnpj:     "49073738000178",
			expected: nil,
		},
		{
			name:     "Valid CNPJ with formatting",
			cnpj:     "02.550.635/0001-98",
			expected: nil,
		},
		{
			name:     "Invalid CNPJ - Wrong check digits",
			cnpj:     "40743953000161",
			expected: ErrInvalidDocument,
		},
		{
			name:     "Invalid CNPJ - Repeated digits",
			cnpj:     "11111111111111",
			expected: ErrRepeatedDigits,
		},
		{
			name:     "Invalid CNPJ - Incorrect length",
			cnpj:     "4074395300016",
			expected: ErrInvalidDocument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCNPJ(tc.cnpj)
			if err != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, err)
			}
		})
	}
}

// Benchmark para ValidateCPF
func BenchmarkValidateCPF(b *testing.B) {
	cpf := "52998224725"
	for i := 0; i < b.N; i++ {
		ValidateCPF(cpf)
	}
}

// Benchmark para ValidateCNPJ
func BenchmarkValidateCNPJ(b *testing.B) {
	cnpj := "40743953000160"
	for i := 0; i < b.N; i++ {
		ValidateCNPJ(cnpj)
	}
}
