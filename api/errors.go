package api

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type ApiError struct {
	Type    string
	Message string
}

func (k ApiError) Error() string {
	return fmt.Sprintf(`ApiError: {Type: "%s", Message: "%s"}`, k.Type, k.Message)
}

type KeyNotFoundError struct {
	Key string
}

func (k KeyNotFoundError) Error() string {
	return fmt.Sprintf(`couldn't find the error message with the key "%s"`, k.Key)
}

type PoorlyFormattedCSV struct{}

func (k PoorlyFormattedCSV) Error() string {
	return fmt.Sprintf(`This CSV is poorly formatted, it should have at least 2 columns on every row`)
}

var UndefinedApiError = ApiError{Type: "UndefinedApiError", Message: "This error has no description setup"}

const (
	ApiErrorWrongParamType          = "WrongParamType"
	ApiRestrictedArea               = "RestrictedArea"
	ApiWrongCredentials             = "WrongCredentials"
	ApiPasswordConfirmationMismatch = "PasswordConfirmationMismatch"
)

func getErrorMessages() ([][]string, error) {
	file, err := os.Open("settings/error_messages.csv")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	errorMessages, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return errorMessages, nil
}

// Reads the folder settings/error_messages.csv
func GetApiError(key string) ApiError {
	errorMessages, err := getErrorMessages()
	if err != nil {
		log.Default().Printf(`Couldn't Get Api Error key: %s got error: %v`, key, err)
		return UndefinedApiError
	}

	for i := 0; i < len(errorMessages); i++ {
		if key == errorMessages[i][0] {
			if len(errorMessages[i]) < 2 {
				log.Default().Printf(`Couldn't Get Api Error key: %s got error: %v`, key, &PoorlyFormattedCSV{})
				return UndefinedApiError
			}
			return ApiError{Type: key, Message: errorMessages[i][1]}
		}
	}

	log.Default().Printf(`Couldn't Get Api Error key: %s got error: %v`, key, &KeyNotFoundError{})
	return UndefinedApiError
}
