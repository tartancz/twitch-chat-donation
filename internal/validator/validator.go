package validator

import (
	"strings"
	"time"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) Error() string {
	b := strings.Builder{}
	if len(v.FieldErrors) > 0 {
		b.WriteString("Field Errors:\n")
		for k, v := range v.FieldErrors {
			b.WriteString(k + ": " + v + "\n")
		}
	}
	if len(v.NonFieldErrors) > 0 {
		b.WriteString("Non Field Errors:\n")
		for _, v := range v.NonFieldErrors {
			b.WriteString(v + "\n")
		}
	}
	return b.String()
}

func HandleDateRange(validator *Validator, from, to string, fromTime *time.Time, toTime *time.Time) {
	if from == "" {
		*fromTime = time.Time{}
	} else {
		validator.CheckField(ValidAndConvertDateTime(from, time.DateOnly, fromTime), "from", "Invalid from date format. Use 'YYYY-MM-DD' format.")
	}

	if to == "" {
		*toTime = time.Now()
	} else {
		validator.CheckField(ValidAndConvertDateTime(to, time.DateOnly, toTime), "to", "Invalid to date format. Use 'YYYY-MM-DD' format.")
	}
	validator.CheckField(fromTime.Before(*toTime), "from", "From date must be before To date.")
}

func ValidAndConvertDateTime(date string, format string, t *time.Time) bool {
	parsed, err := time.Parse(format, date)
	if err != nil {
		return false
	}
	*t = parsed
	return true
}
