package validator

import (
	"regexp"
)

type InterfaceValidator[T any] interface {
	Validate(*T)
}

type BaseValidator struct {
	Errors map[string]string
}

func NewBaseValidator() *BaseValidator {
	return &BaseValidator{Errors: make(map[string]string)}
}

func (v *BaseValidator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (v *BaseValidator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *BaseValidator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *BaseValidator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *BaseValidator) In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

//func Matches(value string, rx *regexp.Regexp) bool {
//	return rx.MatchString(value)
//}
//
//func Unique(values []string) bool {
//	uniqueValues := make(map[string]bool)
//	for _, value := range values {
//		uniqueValues[value] = true
//	}
//	return len(values) == len(uniqueValues)
//}
