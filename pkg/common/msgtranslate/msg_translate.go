package msgtranslate

import (
	"fmt"

	"backend-service/pkg/common/logger"

	"golang.org/x/text/language"

	"github.com/ghodss/yaml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	singleton  *Translator
	langTagMap = map[string]language.Tag{
		"vi": language.Vietnamese,
		"en": language.English,
	}
)

const defaultLang = "en"

type Translator struct {
	translators map[string]*i18n.Localizer
}

// Translator use i18n to translate message by language code.
// Currently support [en, vn].
func Translate(translationKey string, lang *string, data map[string]any) string {
	langKey := defaultLang
	if lang != nil && *lang != "" {
		langKey = *lang
	}

	translator := singleton.translators[langKey]
	if translator == nil {
		return translationKey
	}

	res, err := translator.Localize(
		&i18n.LocalizeConfig{
			MessageID:    translationKey,
			TemplateData: data,
		},
	)
	if err != nil {
		logger.Errorw(err.Error(), "translationKey", translationKey, "data", data)

		return translationKey
	}

	return res
}

// Init is Translator constructor.
func Init() *Translator {
	translators := make(map[string]*i18n.Localizer)

	for lang, tag := range langTagMap {
		bundle := i18n.NewBundle(tag)
		bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

		if _, err := bundle.LoadMessageFile(fmt.Sprintf("translation.%s.yaml", lang)); err != nil {
			logger.Errorw("fail to load translation keys", "err", err, "lang", lang)

			return nil
		}

		translators[lang] = i18n.NewLocalizer(bundle, lang)
	}

	singleton = &Translator{translators: translators}

	return singleton
}
