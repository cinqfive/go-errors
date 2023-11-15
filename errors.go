package errors

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tuikart8/vertical-sphere/utils"
)

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type WebError struct {
	Status int    // the HTTP status code applicable to this problem, expressed as a string value.
	Code   string // an application-specific error code, expressed as a string value.
	Title  string //a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization.
	Detail string // a human-readable explanation specific to this occurrence of the problem. Like title,
	Type   string
}

type WebFieldsError struct {
	Type   string       `json:"type"`
	Status int          `json:"status"` // the HTTP status code applicable to this problem, expressed as a string value.
	Code   string       `json:"code"`   // an application-specific error code, expressed as a string value.
	Title  string       `json:"title"`  // short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization.
	Detail string       `json:"detail"` // a human-readable explanation specific to this occurrence of the problem. Like title,
	Errors []FieldError `json:"errors"`
}

type ErrorDescription struct {
	Code        string
	Title       string
	Description string
}

type ErrorPageData struct {
	Error WebError
}

var errorsMap map[string]ErrorDescription

func SendError(status int, code string, w http.ResponseWriter) {
	const errorType = "fb.entities.WebError"
	var errorDescription = errorsMap[code]
	SendPreparedError(
		WebError{
			Status: status,
			Code:   code,
			Title:  errorDescription.Title,
			Type:   errorType,
			Detail: errorDescription.Description,
		},
		w,
	)
}

func RenderError(status int, code string, w http.ResponseWriter) {
	errorDescription := errorsMap[code]
	data := ErrorPageData{
		Error: WebError{
			Status: status,
			Code:   code,
			Title:  errorDescription.Title,
			Detail: errorDescription.Description,
			Type:   "github.com/cinqfive/go-errors/GeneralError",
		},
	}

	if err := utils.RenderTemplate("error", data, w); err != nil {
		SendError(status, code, w)
	}
}

func SendPreparedError(webError interface{}, w http.ResponseWriter) {
	sendResponse(webError, webError.(WebError).Status, w)
}

func sendResponse(data interface{}, statusCode int, w http.ResponseWriter) {
	jsondata, _ := json.Marshal(data)
	w.WriteHeader(statusCode)
	w.Write(jsondata)
}

func LoadErrorDescriptions() {
	errorsMap = make(map[string]ErrorDescription)
	descriptions := readErrorFile()
	for _, descObj := range descriptions {
		errorsMap[descObj.Code] = descObj
	}
}

func readErrorFile() []ErrorDescription {
	var descriptions []ErrorDescription
	jsonFile, err := os.Open("errors.json")
	if err != nil {
		log.Fatalf("error reading file json: %v", err)
	}

	defer jsonFile.Close()
	bytesValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(bytesValue, &descriptions)

	return descriptions
}

func sendFieldErrors(
	status int,
	errorCode string,
	fieldErrors []FieldError,
	w http.ResponseWriter,
) {
	const errorType = "fb.entities.WebError"
	var errorDescription = errorsMap[errorCode]

	sendPreparedFieldsError(
		WebFieldsError{
			Status: status,
			Code:   errorCode,
			Title:  errorDescription.Description,
			Type:   errorType,
			Detail: errorDescription.Description,
			Errors: fieldErrors,
		},
		w,
	)
}

func sendPreparedFieldsError(WebError WebFieldsError, w http.ResponseWriter) {
	sendResponse(WebError, WebError.Status, w)
}
