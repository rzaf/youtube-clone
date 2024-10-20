package helpers

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/rzaf/youtube-clone/file/models"
	"log"
	"os"
	"strconv"
)

const (
	OneMB               = 1000000
	UserMaxUploadIn24MB = 10
	UserMaxUploadIn24B  = UserMaxUploadIn24MB * OneMB
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

func ValidateInt(str string, param string) int {
	if str == "" {
		return 0
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		panic(ValidationFieldError{param, "should be a integer"})
	}
	return n
}

func ValidatePositiveInt(str string, param string) int {
	if str == "" {
		return 0
	}
	n, err := strconv.Atoi(str)
	if err != nil || n < 1 {
		panic(ValidationFieldError{param, "should be a positive integer"})
	}
	return n
}

func ValidateVideoType(filetype string) {
	switch filetype {
	case "video/mp4", "video/mpeg", "video/x-msvideo", "video/webm":
		return
	default:
		panic(ServerError{Message: fmt.Sprintf("invalid video type:'%s' . supported types:mp4,mpeg,webm,x-msvideo", filetype), Status: 400})
	}
}

func ValidateImageType(filetype string) {
	switch filetype {
	case "image/jpeg", "image/jpg", "image/gif", "image/png":
		return
	default:
		panic(ServerError{Message: fmt.Sprintf("invalid image type:'%s' . supported types:jpeg,jpg,gif,png", filetype), Status: 400})
	}
}

func ValidateMusicType(filetype string) {
	switch filetype {
	case "audio/mpeg", "audio/x-ms-wma", "audio/vnd.rn-realaudio", "audio/x-wav":
		return
	default:
		panic(ServerError{Message: fmt.Sprintf("invalid audio type:'%s' . supported types:jpeg,jpg,gif,png", filetype), Status: 400})
	}
}

func ValidateUrl(url string) {
	if url == "" {
		panic(NewServerError("Invalid Url", 400))
	}
	if len(url) != 16 {
		panic(NewServerError("Invalid Url", 400))
	}
}

func ValidateVideoUrl(url string) {
	if url == "" {
		panic(NewServerError("Invalid Url", 400))
	}
	urlLen := len(url)
	if urlLen == 16 {
		return
	}
	if urlLen == 23 {
		if url[urlLen-3:] != ".ts" {
			panic(NewServerError("Invalid Segment Url", 400))
		}
		return
	}
	panic(NewServerError("Invalid Url", 400))
}

func CheckUserUploadBandwidth(size int64, userId int64) {
	res, err := models.GetUserUploadSizeIn24(userId)
	if err != nil {
		panic(err)
	}
	fmt.Println("user total upload size:", res)
	fmt.Println(res+size, UserMaxUploadIn24B, UserMaxUploadIn24MB)
	if res+size > UserMaxUploadIn24B {
		panic(NewServerError("You reached maximum upload size of "+strconv.Itoa(UserMaxUploadIn24MB)+"MB in 24h", 400))
	}
}
func FatalIfEmptyVar(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is not set !\n", key)
	}
	return v
}
