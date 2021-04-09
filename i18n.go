package zelduh

import "errors"

type LocaleMessagesMap map[string]string

// just a stub for now since English is the only language supported at this time
var localeMessagesByLanguage = map[string]LocaleMessagesMap{
	"en": {
		"gameTitle":             "Zelduh",
		"pauseScreenMessage":    "Paused",
		"gameOverScreenMessage": "Game Over",
	},
	"es": {
		"gameTitle":             "Zelduh",
		"pauseScreenMessage":    "Paused",
		"gameOverScreenMessage": "Game Over",
	},
}

// GetLocaleMessageMapByLanguage returns a map of message IDs to translation strings by language
func GetLocaleMessageMapByLanguage(language string) (LocaleMessagesMap, error) {
	if language != "en" && language != "es" {
		return nil, errors.New("language not supported")
	}
	return localeMessagesByLanguage[language], nil
}
