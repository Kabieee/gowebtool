package common

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTran "github.com/go-playground/validator/v10/translations/en"
	_ "github.com/go-playground/validator/v10/translations/zh"
)

var (
	MyTran ut.Translator
)

func InitTranslator() {
	if va, ok := binding.Validator.Engine().(*validator.Validate); ok {
		loc := en.New()
		u := ut.New(loc, loc)
		translator, _ := u.GetTranslator(loc.Locale())
		_ = enTran.RegisterDefaultTranslations(va, translator)
		va.RegisterTagNameFunc(func(field reflect.StructField) string {
			return fmt.Sprintf("[%s]", field.Name)
		})
		MyTran = translator
	}
}
