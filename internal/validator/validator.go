package validator

import "slices"

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if ok {
		return
	}

	v.AddError(key, message)
}

func PermittedValue[T comparable](value T, permittedValiues ...T) bool {
	return slices.Contains(permittedValiues, value)
}
