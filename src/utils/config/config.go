package config

import (
	"encoding/json"
	"os"

	"github.com/irdaislakhuafa/go-sdk/caches"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/db"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/files"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/smtp"
	"github.com/irdaislakhuafa/go-sdk/storage"
)

type AppMode string

const (
	APP_MODE_DEV  = "dev"
	APP_MODE_PROD = "prod"
)

type (
	Meta struct {
		Title       string
		Description string
		Version     string
		Protocol    string
		Port        string
		Host        string
		BasePath    string
	}

	Fiber struct {
		Port           string
		TimeoutSeconds int
		Mode           string
		Cors           struct {
			PreferDefault bool
			Origins       []string
			Methods       []string
			Headers       []string
			Cookies       bool
		}
	}

	DB struct {
		Master db.Config
	}

	Token struct {
		ExpirationMinutes int
		Secret            string
	}

	Contacts struct {
		Name  string
		Email string
	}

	Config struct {
		Meta     Meta
		Fiber    Fiber
		SMTP     smtp.Config
		Log      log.Config
		DB       DB
		Storage  storage.Config
		Cache    caches.Config
		Token    Token
		Contacts Contacts
	}
)

func ReadFileJSON(pathToFile string) (Config, error) {
	if !files.IsExist(pathToFile) {
		return Config{}, errors.NewWithCode(codes.CodeStorageNoFile, "file '%s' not found", pathToFile)
	}

	fileBytes, err := os.ReadFile(pathToFile)
	if err != nil {
		return Config{}, errors.NewWithCode(codes.CodeStorageNoFile, "cannot read file '%s': %s", pathToFile, err.Error())
	}

	result := Config{}
	if err := json.Unmarshal(fileBytes, &result); err != nil {
		return Config{}, errors.NewWithCode(codes.CodeJSONUnmarshalError, "cannot parse '%s': %s", pathToFile, err.Error())
	}

	return result, nil
}
