package helpers

import (
	"github.com/rzaf/youtube-clone/notification/pbs/notificationHelperPb"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
)

var validate *validator.Validate = validator.New()

type ValidationFieldError struct {
	Param   string
	Message string
}

type ValidationFieldErrors []ValidationFieldError

func (v *ValidationFieldErrors) ErrorMessage() any {
	return map[string]any{
		"error": *v,
	}
}

func (v *ValidationFieldError) ErrorMessage() any {
	return map[string]any{
		"error": *v,
	}
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "max":
		return "max is " + fe.Param()
	case "min":
		return "min is " + fe.Param()
	case "url":
		return "Invalid url" + fe.Param()
	}
	return ""
	// return fe.Error() // default error
}

func ValidateVar(field interface{}, paramName string, tag string) {
	if err := validate.Var(field, tag); err != nil {
		ve, _ := err.(validator.ValidationErrors)
		// return &ValidationFieldError{paramName, msgForTag(ve[0])}
		panic(ValidationFieldError{paramName, msgForTag(ve[0])})
	}
}

func ValidateAllowedParams(m map[string]any, allowedParams ...string) {
	paramMap := make(map[string]bool, 5)
	for _, p := range allowedParams {
		paramMap[p] = true
	}
	for p := range m {
		if _, ok := paramMap[p]; !ok {
			panic(ValidationFieldError{p, "is not a valid parameter! valid params are: " + strings.Join(allowedParams, ",")})
		}
	}
}

func ValidateStr(strInteface any, param string, defaultValue string) string {
	switch str := strInteface.(type) {
	case string:
		if str == "" {
			panic(ValidationFieldError{param, "is empty"})
			// return defaultValue
		}
		return str
	case nil:
		return defaultValue
	}
	panic(ValidationFieldError{param, "should be string"})
}

// ensures string not be empty
func ValidateRequiredStr(strInteface any, param string) string {
	switch str := strInteface.(type) {
	case string:
		if str == "" {
			panic(ValidationFieldError{param, "is empty"})
		}
		return str
	case nil:
		panic(ValidationFieldError{param, "is required"})
	}
	panic(ValidationFieldError{param, "should be string"})
}

func ValidateInt(intInterface any, param string, defaultValue int) int {
	switch v := intInterface.(type) {
	case string:
		if v == "" {
			panic(ValidationFieldError{param, "is empty"})
			// return defaultValue
		}
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	case float64:
		return int(v)
	case nil:
		return defaultValue
	}
	panic(ValidationFieldError{param, "should be a integer"})
}

func ValidateRequiredInt(intInterface any, param string) int {
	switch v := intInterface.(type) {
	case string:
		if v == "" {
			panic(ValidationFieldError{param, "is empty"})
		}
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	case float64:
		return int(v)
	case nil:
		panic(ValidationFieldError{param, "is required"})
	}
	panic(ValidationFieldError{param, "should be a integer"})
}

func ValidatePositiveInt(intInterface any, param string, defaultValue int) int {
	switch v := intInterface.(type) {
	case string:
		if v == "" {
			panic(ValidationFieldError{param, "is empty"})
			// return defaultValue
		}
		n, err := strconv.Atoi(v)
		if err == nil && n > 0 {
			return n
		}
	case float64:
		n := int(v)
		if n > 0 {
			return n
		}
	case nil:
		return defaultValue
	}
	panic(ValidationFieldError{param, "should be a positive integer"})
}

func ValidateRequiredPositiveInt(intInterface any, param string) int {
	switch v := intInterface.(type) {
	case string:
		if v == "" {
			panic(ValidationFieldError{param, "is empty"})
		}
		n, err := strconv.Atoi(v)
		if err == nil && n > 0 {
			return n
		}
	case float64:
		n := int(v)
		if n > 0 {
			return n
		}
	case nil:
		panic(ValidationFieldError{param, "is required"})
	}
	panic(ValidationFieldError{param, "should be a positive integer"})
}

func ValidateBool(boolInterface any, param string) bool {
	switch v := boolInterface.(type) {
	case string:
		if v == "" {
			panic(ValidationFieldError{param, "is empty"})
		}
		n, err := strconv.ParseBool(v)
		if err == nil {
			return n
		}
	case float64:
		if v == 0 {
			return false
		}
		if v == 1 {
			return true
		}
	case bool:
		return v
	case nil:
		panic(ValidationFieldError{param, "is required"})
	}
	panic(ValidationFieldError{param, `should be boolean (1,0,true,false,"1","0","true","false")`})
}

func ValidateNotificationSortTypes(t string) notificationHelperPb.SortType {
	switch t {
	case "newest":
		return notificationHelperPb.SortType(notificationHelperPb.SortType_Newest)
	case "oldest":
		return notificationHelperPb.SortType(notificationHelperPb.SortType_Oldest)
	}
	panic(ValidationFieldError{"type", "incorrect sort type:`" + t + "`"})
}

func ValidateNotificationSeenTypes(t string) notificationHelperPb.SeenType {
	switch t {
	case "any":
		return notificationHelperPb.SeenType(notificationHelperPb.SeenType_Any)
	case "seen":
		return notificationHelperPb.SeenType(notificationHelperPb.SeenType_Seen)
	case "not-seen":
		return notificationHelperPb.SeenType(notificationHelperPb.SeenType_NotSeen)
	}
	panic(ValidationFieldError{"type", "incorrect notification type:`" + t + "`"})
}

func FatalIfEmptyVar(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is not set !\n", key)
	}
	return v
}
