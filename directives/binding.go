package directives

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)
}

func Binding(ctx context.Context, obj interface{}, next graphql.Resolver, constraint string, trim *bool) (res interface{}, err error) {
	val, err := next(ctx)
	if err != nil {
		return val, err
	}
	fieldName := *graphql.GetPathContext(ctx).Field

	if trim != nil && *trim {
		tempVal, ok := val.(string)
		if !ok {
			tempVal, ok := val.(*string)
			if !ok {
				return nil, fmt.Errorf("[trim] failed, %s is not a string", fieldName)
			}
			tempStr := strings.Trim(*tempVal, " ")
			val = &tempStr
		} else {
			val = strings.Trim(tempVal, " ")
		}
	}

	err = validate.Var(val, constraint)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		transErr := fmt.Errorf("%s%+v", fieldName, validationErrors[0].Translate(trans))
		return val, transErr
	}

	return val, nil
}
