package validateErrors

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (v *ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return ""
	}

	return "validation errors"
}
