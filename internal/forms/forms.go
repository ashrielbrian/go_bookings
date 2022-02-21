package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

// Valid checks if the form is valid
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initialises a new Form based on the form fields submitted by the user
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank!")
		}
	}
}

// Has checks if a field exists in the request body
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	return x != ""
}

// MinLength checks for a minimum character length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for invalid email address
func (f *Form) IsEmail(field string) bool {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
		return false
	}
	return true
}
