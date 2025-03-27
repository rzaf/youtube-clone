package helpers

import (
	"github.com/rzaf/youtube-clone/database/pbs/helper"
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

func ValidateUsersSortTypes(t string) helper.SortType {
	switch t {
	case "newest":
		return helper.SortType(helper.SortType_Newest)
	case "oldest":
		return helper.SortType(helper.SortType_Oldest)
	case "most-viewed":
		return helper.SortType(helper.SortType_MostViewed)
	case "least-viewed":
		return helper.SortType(helper.SortType_LeastViewed)
	case "most-subbed":
		return helper.SortType(helper.SortType_MostSubscribers)
	case "least-subbed":
		return helper.SortType(helper.SortType_LeastSubscribers)
	}
	panic(ValidationFieldError{"type", "incorrect sort type:`" + t + "`"})
}

func ValidateMediasSortTypes(t string) helper.SortType {
	switch t {
	case "newest":
		return helper.SortType(helper.SortType_Newest)
	case "oldest":
		return helper.SortType(helper.SortType_Oldest)
	case "most-viewed":
		return helper.SortType(helper.SortType_MostViewed)
	case "least-viewed":
		return helper.SortType(helper.SortType_LeastViewed)
	}
	panic(ValidationFieldError{"type", "incorrect sort type:`" + t + "`"})
}

func ValidateCommentsSortTypes(t string) helper.SortType {
	switch t {
	case "newest":
		return helper.SortType(helper.SortType_Newest)
	case "oldest":
		return helper.SortType(helper.SortType_Oldest)
	case "most-liked":
		return helper.SortType(helper.SortType_MostLiked)
	case "least-liked":
		return helper.SortType(helper.SortType_LeastLiked)
	case "most-disliked":
		return helper.SortType(helper.SortType_MostDisiked)
	case "least-disliked":
		return helper.SortType(helper.SortType_LeastDisliked)
	case "most-replied":
		return helper.SortType(helper.SortType_MostReplied)
	case "least-replied":
		return helper.SortType(helper.SortType_LeastReplied)
	}
	panic(ValidationFieldError{"type", "incorrect sort type:`" + t + "`"})
}

func ValidatePlaylistsSortTypes(t string) helper.SortType {
	switch t {
	case "newest":
		return helper.SortType(helper.SortType_Newest)
	case "oldest":
		return helper.SortType(helper.SortType_Oldest)
	case "most-viewed":
		return helper.SortType(helper.SortType_MostViewed)
	case "least-viewed":
		return helper.SortType(helper.SortType_LeastViewed)
	}
	panic(ValidationFieldError{"type", "incorrect sort type:`" + t + "`"})
}

func ValidateMediaType(t string) helper.MediaType {
	switch t {
	case "video":
		return helper.MediaType(helper.MediaType_VIDEO)
	case "music":
		return helper.MediaType(helper.MediaType_MUSIC)
	case "photo":
		return helper.MediaType(helper.MediaType_PHOTO)
	}
	panic(ValidationFieldError{"type", "incorrect media type:`" + t + "`"})
}

// only in queries that accept all media types
func ValidateAllMediaType(t string) helper.MediaType {
	switch t {
	case "video":
		return helper.MediaType(helper.MediaType_VIDEO)
	case "music":
		return helper.MediaType(helper.MediaType_MUSIC)
	case "photo":
		return helper.MediaType(helper.MediaType_PHOTO)
	case "all", "any":
		return helper.MediaType(helper.MediaType_ALL)
	}
	panic(ValidationFieldError{"type", "incorrect media type:`" + t + "`"})
}

func FatalIfEmptyVar(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is not set !\n", key)
	}
	return v
}
