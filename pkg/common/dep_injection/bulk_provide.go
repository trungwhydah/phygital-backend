package depinjection

import (
	"go.uber.org/fx"
)

// BulkProvide put all constructor defined to fx pool and use everywhere in this project.
func BulkProvide(providingFunctions []any, group string) fx.Option {
	var fxOptions []fx.Option

	for _, function := range providingFunctions {
		fxOptions = append(
			fxOptions,
			fx.Provide(
				fx.Annotated{
					Target: function,
					Group:  group,
				},
			),
		)
	}

	return fx.Options(fxOptions...)
}
