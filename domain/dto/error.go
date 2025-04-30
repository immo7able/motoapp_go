package dto

type Error struct {
	Code        int          `json:"code"`
	Message     string       `json:"message"`
	FieldErrors []FieldError `json:"fields"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
