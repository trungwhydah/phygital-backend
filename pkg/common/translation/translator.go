package translation

import (
	"backend-service/pkg/common/logger"
	structtraversal2 "backend-service/pkg/marketplace/structtraversal"

	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
)

type Translatable interface {
	GetTranslations() Translations
}

type Translations map[string]map[string]string

// TranslateCollection is for translating all elements in a slice.
func TranslateCollection(coll any, lang string) {
	structtraversal2.TraverseSlice(coll, translateFieldCallback(lang))
}

// Translate is for translating content of an object.
func Translate(obj any, lang string) {
	structtraversal2.TraverseObject(obj, translateFieldCallback(lang))
}

func translateFieldCallback(lang string) func(args ...any) {
	return func(args ...any) {
		if len(args) == 0 {
			return
		}

		field := args[0]
		fieldVal, ok := field.(Translatable)

		if !ok {
			return
		}

		translation, ok := fieldVal.GetTranslations()[lang]
		if !ok {
			return
		}

		err := mapstructure.Decode(translation, fieldVal)
		if err != nil {
			logger.Errorw(
				err.Error(),
				"field", fieldVal,
				"translation", translation,
			)

			return
		}

		if err := copier.CopyWithOption(
			fieldVal,
			translation,
			copier.Option{IgnoreEmpty: true},
		); err != nil {
			logger.Errorw(
				err.Error(),
				"field", fieldVal,
				"translation", translation,
			)
		}
	}
}
