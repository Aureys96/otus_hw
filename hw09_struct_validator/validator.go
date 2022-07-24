package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	errWrongType               = errors.New("value for validation is not a struct")
	errStringWrongLength       = errors.New("string has invalid length")
	errStringDoesntMatchRegexp = errors.New("string doesn't match pattern")
	errUnexpectedValue         = errors.New("unexpected value")
	errInvalidMinIntValue      = errors.New("int value is less than min")
	errInvalidMaxIntValue      = errors.New("int value is greater than max")
)

type ProgrammaticError struct {
	root interface{}
}

func (err ProgrammaticError) Error() string {
	root := err.root
	switch v := root.(type) {
	case error:
		return fmt.Sprintf("Received programmatic error: %v\n", v.Error())
	case string:
		return fmt.Sprintf("Received programmatic error: %v\n", v)
	default:
		return "Received unexpected exception"
	}
}

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) hasError() bool {
	return v.Err != nil
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for _, e := range v {
		builder.WriteString("field: ")
		builder.WriteString(e.Field)
		builder.WriteString(", error:  ")
		builder.WriteString(e.Err.Error())
	}
	return builder.String()
}

func Validate(v interface{}) (resultErr error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
			resultErr = ProgrammaticError{root: err}
		}
	}()

	value := reflect.ValueOf(v)
	valueType := value.Type()

	if value.Kind() != reflect.Struct {
		return ProgrammaticError{errWrongType}
	}

	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < valueType.NumField(); i++ {
		fieldValueType := valueType.Field(i)
		strTags := fieldValueType.Tag.Get("validate")
		if len(strTags) == 0 {
			continue
		}
		fieldValue := value.Field(i)

		tags := strings.Split(strTags, "|")

		var validationError ValidationError
		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.String:
			log.Println("received a string")
			if validationError = validateString(tags, fieldValueType.Name, fieldValue.String()); validationError.hasError() {
				validationErrors = append(validationErrors, validationError)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			log.Println("received an int")
			if validationError = validateInteger(tags, fieldValueType.Name, fieldValue.Int()); validationError.hasError() {
				validationErrors = append(validationErrors, validationError)
			}
		case reflect.Slice:
			validationErrors = handleSlices(fieldValue, tags, fieldValueType, validationErrors)
		default:
			log.Printf("received unsupported type: %v\n", fieldValue.Type().String())
		}
	}
	if len(validationErrors) > 0 {
		resultErr = validationErrors
	}
	return resultErr
}

func handleSlices(fieldValue reflect.Value, tags []string, fieldValueType reflect.StructField,
	validationErrors ValidationErrors,
) ValidationErrors {
	var validationError ValidationError
	switch fieldValue.Type().String() {
	case "[]string":
		log.Println("received a string slice")
		for i := 0; i < fieldValue.Len(); i++ {
			validationError = validateString(tags, fieldValueType.Name, fieldValue.Index(i).String())
			if validationError.hasError() {
				validationErrors = append(validationErrors, validationError)
			}
		}
	case "[]int", "[]int8", "[]int16", "[]in32", "[]int64":
		log.Println("received an int slice")
		for i := 0; i < fieldValue.Len(); i++ {
			validationError = validateInteger(tags, fieldValueType.Name, fieldValue.Index(i).Int())
			if validationError.hasError() {
				validationErrors = append(validationErrors, validationError)
			}
		}
	}
	return validationErrors
}

func validateString(tags []string, field string, value string) ValidationError {
	for _, tag := range tags {
		t := strings.Split(tag, ":")
		var err error
		switch t[0] {
		case "len":
			if err = validateStringLength(value, t[1]); err != nil {
				return ValidationError{
					Field: field,
					Err:   err,
				}
			}
		case "regexp":
			if err = validateStringPattern(t, value); err != nil {
				return ValidationError{
					Field: field,
					Err:   err,
				}
			}
		case "in":
			if err = validateStringSet(t, value); err != nil {
				return ValidationError{
					Field: field,
					Err:   err,
				}
			}
		default:
			log.Panicln("Received unknown validation token")
		}
	}
	return ValidationError{}
}

func validateStringSet(t []string, value string) error {
	strs := strings.Split(t[1], ",")
	set := make(Set)
	for _, str := range strs {
		set.add(str)
	}
	if set.contains(value) {
		return nil
	}
	return errUnexpectedValue
}

func validateStringPattern(t []string, value string) error {
	rg, err := regexp.Compile(t[1])
	if err != nil {
		log.Panicf("Provided regexp is invalid: %v\n", t[1])
	}
	if matched := rg.MatchString(value); !matched {
		return errStringDoesntMatchRegexp
	}
	return nil
}

func validateStringLength(value string, sLength string) error {
	length, err := strconv.Atoi(sLength)
	if err != nil {
		log.Panicln("Couldn't convert string length to integer")
	}

	if len(value) != length {
		return errStringWrongLength
	}

	return nil
}

func validateInteger(tags []string, field string, value int64) ValidationError {
	for _, tag := range tags {
		t := strings.Split(tag, ":")
		var err error
		switch t[0] {
		case "min":
			if err = validateIntegerMinMax(value, t[1], false); err != nil {
				return ValidationError{
					Field: field,
					Err:   err,
				}
			}
		case "max":
			if err = validateIntegerMinMax(value, t[1], true); err != nil {
				return ValidationError{
					Field: field,
					Err:   err,
				}
			}
		case "in":
			if err = validateIntSet(t[1], value); err != nil {
				return ValidationError{
					Field: field,
					Err:   err,
				}
			}
		default:
			log.Panicln("Received unknown validation token")
		}
	}
	return ValidationError{}
}

func validateIntegerMinMax(value int64, thresholdValue string, isMax bool) error {
	cond, err := strconv.Atoi(thresholdValue)
	if err != nil {
		log.Panicf("Couldn't convert threshold validation value to int: %v\n", thresholdValue)
	}
	if isMax {
		if int(value) > cond {
			return errInvalidMaxIntValue
		}
	} else {
		if int(value) < cond {
			return errInvalidMinIntValue
		}
	}
	return nil
}

func validateIntSet(t string, value int64) error {
	strs := strings.Split(t, ",")
	set := make(Set)
	for _, str := range strs {
		setValue, err := strconv.Atoi(str)
		if err != nil {
			log.Panicf("Couldn't convert set value to int: %v\n", setValue)
		}
		set.add(int64(setValue))
	}
	if set.contains(value) {
		return nil
	}
	return errUnexpectedValue
}

type Set map[interface{}]interface{}

func (s *Set) add(in interface{}) {
	(*s)[in] = struct{}{}
}

func (s *Set) contains(in interface{}) bool {
	_, contains := (*s)[in]
	return contains
}
