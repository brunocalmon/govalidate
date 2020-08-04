package govalidate

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	emailRegex       = "^(([^<>()\\[\\]\\.,;:\\s@\\\"]+(\\.[^<>()\\[\\]\\.,;:\\s@\\\"]+)*)|(\\\".+\\\"))@(([^<>()[\\]\\.,;:\\s@\\\"]+\\.)+[^<>()[\\]\\.,;:\\s@\\\"]{2,})$" // email pattern
	dateOfBirthRegex = "^(0[1-9]|[12][0-9]|3[01])[\\.](0[1-9]|1[012])[\\.](19|20)\\d\\d$"                                                                                  // german date pattern

	specialCharacters = "(.*[^A-Za-z0-9].*)" // has special char
	upperCase         = "(.*[A-Z].*)"        // has uppercase char
	lowerCase         = "(.*[a-z].*)"        // has lowercase char
	digit             = "(.*\\d.*)"          // has digit
	whitespace        = "(.*[\\s].*)"        // has whitespace
)

type validator interface {
	validate([]string, interface{}) error
}

type defaultValidator struct{}
type stringValidator struct{}
type numberValidator struct{}

func (v defaultValidator) validate(conditions []string, val interface{}) error {
	for _, condition := range conditions {
		switch condition {
		case "required":
			if val == nil {
				return errors.New("field [%s] is required, but is missing")
			}
		}
	}
	return nil
}

func (v stringValidator) validate(conditions []string, val interface{}) error {
	for _, condition := range conditions {
		switch condition {
		case "required":
			if val == "" {
				return errors.New("field [%s] is required, but is missing")
			}
		case "email":
			if val != "" {
				matched, err := regexp.MatchString(emailRegex, val.(string))
				if !matched || err != nil {
					return errors.New("field [%s] should be a email")
				}
			}
		case "password":
			if val != "" {
				password := val.(string)
				matched, err := PasswordValidator(password)
				if !matched || err != nil {
					return errors.New("field [%s] should be a valid password between 8 to 20 characters which contain at least one lowercase letter, one uppercase letter, one numeric digit, and one special character")
				}
			}
		case "gender":
			if val != "" && !(strings.EqualFold("M", val.(string)) || strings.EqualFold("F", val.(string))) {
				return errors.New("field [%s] only accepts 'M' or 'F'")
			}
		case "date":
			if val != "" {
				matched, err := regexp.MatchString(dateOfBirthRegex, val.(string))
				if !matched || err != nil {
					return errors.New("field [%s] should be in the format [dd.mm.yyyy]")
				}
			}
		default:
			err := validateMinMaxLength(condition, val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v numberValidator) validate(conditions []string, val interface{}) error {
	for _, condition := range conditions {
		switch {
		case strings.Contains(condition, "min=") || strings.Contains(condition, "max="):
			constraints := strings.Split(condition, "=")
			fieldVal, ok := val.(int)

			offset, atoiErr := strconv.Atoi(constraints[1])
			if len(constraints) < 2 || !ok || atoiErr != nil {
				return errors.New("field [%s] validations malformed")
			}

			if constraints[0] == "min" {
				if fieldVal < offset {
					return errors.New(fmt.Sprintf("the field [%s] value [%d] has range shorter than required [%d]", "%s", fieldVal, offset))
				}
			}
			if constraints[0] == "max" {
				if fieldVal > offset {
					return errors.New(fmt.Sprintf("the field [%s] value [%d] has range higher than required [%d]", "%s", fieldVal, offset))
				}
			}
		}
	}
	return nil
}

func getValidatorByType(fieldType string) validator {
	switch fieldType {
	case "string":
		return stringValidator{}
	case "int":
		return numberValidator{}
	}
	return defaultValidator{}
}

func BodyRequest(d interface{}) []string {
	var fieldErrors []string
	v := reflect.Indirect(reflect.ValueOf(d))

	for i := 0; i < v.NumField(); i++ {
		validate := v.Type().Field(i).Tag.Get("validate")
		conditions := make([]string, 0)
		if validate != "" {
			conditions = strings.Split(validate, ",")
		}

		if len(conditions) < 1 {
			continue
		}

		field := v.Field(i)
		fieldType := field.Type().Name()
		fieldValue := field.Interface()
		fieldName := v.Type().Field(i).Name

		validator := getValidatorByType(fieldType)
		err := validator.validate(conditions, fieldValue)
		if err != nil {
			fieldErrors = append(fieldErrors, fmt.Sprintf(err.Error(), fieldName))
		}
	}

	if fieldErrors != nil {
		return fieldErrors
	}
	return nil
}

func validateMinMaxLength(condition string, val interface{}) error {
	if strings.Contains(condition, "min=") || strings.Contains(condition, "max=") {
		constraints := strings.Split(condition, "=")
		fieldVal, ok := val.(string)

		offset, atoiErr := strconv.Atoi(constraints[1])
		if len(constraints) < 2 || !ok || atoiErr != nil {
			return errors.New("field [%s] validations malformed")
		}

		if constraints[0] == "min" {
			if len(fieldVal) < offset {
				return errors.New(fmt.Sprintf("the field [%s] value [%s] has length shorter than required [%d]", "%s", fieldVal, offset))
			}
		}
		if constraints[0] == "max" {
			if len(fieldVal) > offset {
				return errors.New(fmt.Sprintf("the field [%s] value [%s] has length higher than required [%d]", "%s", fieldVal, offset))
			}
		}
	}
	return nil
}

func PasswordValidator(password string) (bool, error) {
	hasSpecial, errSpecial := regexp.MatchString(specialCharacters, password)
	if errSpecial != nil {
		return false, errSpecial
	}
	hasDigit, errDigit := regexp.MatchString(digit, password)
	if errDigit != nil {
		return false, errDigit
	}
	hasUppercase, errUppercase := regexp.MatchString(upperCase, password)
	if errUppercase != nil {
		return false, errUppercase
	}
	hasLowercase, errLowercase := regexp.MatchString(lowerCase, password)
	if errLowercase != nil {
		return false, errLowercase
	}
	hasWhitespace, errWhitespace := regexp.MatchString(whitespace, password)
	if errWhitespace != nil {
		return false, errWhitespace
	}
	hasCorrectSize := len(password) >= 8 && len(password) <= 20

	log.Printf("Password validation:\n special: [%t]\n digit: [%t]\n uppercase: [%t]\n lowercase: [%t]\n correct size: [%t]\n", hasSpecial, hasDigit, hasUppercase, hasLowercase, hasCorrectSize)
	return hasSpecial && hasDigit && hasUppercase && hasLowercase && hasCorrectSize && !hasWhitespace, nil
}
