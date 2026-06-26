package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

type Language string

const (
	English Language = "en"
	Spanish Language = "es"
)

const (
	KeyMissingCredentials  = "missing_credentials"
	KeyNotFound            = "not_found"
	KeyInvalidInput        = "invalid_input"
	KeyFailedPrecondition  = "failed_precondition"
	KeyInternalServerError = "internal_server_error"
)

type catalog struct {
	Messages map[string]string `json:"messages"`
}

//go:embed locales/*.json
var localeFS embed.FS

var catalogs = mustLoadCatalogs()

func Normalize(value string) Language {
	for _, candidate := range strings.Split(value, ",") {
		current := strings.TrimSpace(candidate)
		if current == "" {
			continue
		}
		if index := strings.IndexByte(current, ';'); index >= 0 {
			current = current[:index]
		}
		current = strings.TrimSpace(strings.ToLower(current))
		if current == "" {
			continue
		}
		if index := strings.IndexAny(current, "-_"); index >= 0 {
			current = current[:index]
		}
		switch Language(current) {
		case Spanish:
			return Spanish
		case English:
			return English
		}
	}
	return English
}

func Message(language Language, key string) string {
	if message, ok := catalogs[language].Messages[key]; ok {
		return message
	}
	return catalogs[English].Messages[key]
}

func mustLoadCatalogs() map[Language]catalog {
	return map[Language]catalog{
		English: mustLoadCatalog(English),
		Spanish: mustLoadCatalog(Spanish),
	}
}

func mustLoadCatalog(language Language) catalog {
	path := fmt.Sprintf("locales/%s.json", language)
	payload, err := localeFS.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("read i18n catalog %s: %v", path, err))
	}
	var current catalog
	if err := json.Unmarshal(payload, &current); err != nil {
		panic(fmt.Sprintf("decode i18n catalog %s: %v", path, err))
	}
	if current.Messages == nil {
		current.Messages = map[string]string{}
	}
	return current
}
