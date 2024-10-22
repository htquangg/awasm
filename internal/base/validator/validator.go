// https://github.com/apache/incubator-answer/blob/main/internal/base/validator/validator.go
package validator

import (
	std_errors "errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/locales"
	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/microcosm-cc/bluemonday"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/i18n"

	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/base/translator"
	"github.com/htquangg/awasm/pkg/logger"
)

type TranslatorLocal struct {
	Lo           locales.Translator
	RegisterFunc func(v *validator.Validate, trans ut.Translator) (err error)
	La           i18n.Language
}

var allLanguageTranslators = []*TranslatorLocal{
	{La: i18n.LanguageEnglish, Lo: english.New(), RegisterFunc: en.RegisterDefaultTranslations},
}

// MyValidator my validator
type MyValidator struct {
	Validate *validator.Validate
	Tran     ut.Translator
	Lang     i18n.Language
}

// FormErrorField indicates the current form error content. which field is error and error message.
type FormErrorField struct {
	ErrorField string `json:"error_field"`
	ErrorMsg   string `json:"error_msg"`
}

// GlobalValidatorMapping is a mapping from validator to translator used
var GlobalValidatorMapping = make(map[i18n.Language]*MyValidator, 0)

func init() {
	for _, t := range allLanguageTranslators {
		tran, val := getTran(t.Lo), createDefaultValidator(t.La)
		if t.RegisterFunc != nil {
			if err := t.RegisterFunc(val, tran); err != nil {
				panic(err)
			}
		}
		GlobalValidatorMapping[t.La] = &MyValidator{Validate: val, Tran: tran, Lang: t.La}
	}
}

func getTran(lo locales.Translator) ut.Translator {
	tran, ok := ut.New(lo, lo).GetTranslator(lo.Locale())
	if !ok {
		panic(fmt.Sprintf("not found translator %s", lo.Locale()))
	}
	return tran
}

func NotBlank(fl validator.FieldLevel) (res bool) {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		trimSpace := strings.TrimSpace(field.String())
		res := len(trimSpace) > 0
		if !res {
			field.SetString(trimSpace)
		}
		return true
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func Sanitizer(fl validator.FieldLevel) (res bool) {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		filter := bluemonday.UGCPolicy()
		content := strings.Replace(filter.Sanitize(field.String()), "&amp;", "&", -1)
		field.SetString(content)
		return true
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func createDefaultValidator(la i18n.Language) *validator.Validate {
	validate := validator.New()
	// _ = validate.RegisterValidation("notblank", validators.NotBlank)
	_ = validate.RegisterValidation("notblank", NotBlank)
	_ = validate.RegisterValidation("sanitizer", Sanitizer)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) (res string) {
		defer func() {
			if len(res) > 0 {
				res = translator.Tr(la, res)
			}
		}()
		if jsonTag := fld.Tag.Get("json"); len(jsonTag) > 0 {
			if jsonTag == "-" {
				return ""
			}
			return jsonTag
		}
		if formTag := fld.Tag.Get("form"); len(formTag) > 0 {
			return formTag
		}
		return fld.Name
	})
	return validate
}

func GetValidatorByLang(lang i18n.Language) *MyValidator {
	if GlobalValidatorMapping[lang] != nil {
		return GlobalValidatorMapping[lang]
	}
	return GlobalValidatorMapping[i18n.DefaultLanguage]
}

// Check /
func (m *MyValidator) Check(value interface{}) (errFields []*FormErrorField, err error) {
	defer func() {
		if len(errFields) == 0 {
			return
		}
		for _, field := range errFields {
			if len(field.ErrorField) == 0 {
				continue
			}
			firstRune := []rune(field.ErrorMsg)[0]
			if !unicode.IsLetter(firstRune) || !unicode.Is(unicode.Latin, firstRune) {
				continue
			}
			upperFirstRune := unicode.ToUpper(firstRune)
			field.ErrorMsg = string(upperFirstRune) + field.ErrorMsg[1:]
			if !strings.HasSuffix(field.ErrorMsg, ".") {
				field.ErrorMsg += "."
			}
		}
	}()
	err = m.Validate.Struct(value)
	if err != nil {
		var valErrors validator.ValidationErrors
		if !std_errors.As(err, &valErrors) {
			logger.Error(err)
			return nil, std_errors.New("validate check exception")
		}

		for _, fieldError := range valErrors {
			errField := &FormErrorField{
				ErrorField: fieldError.Field(),
				ErrorMsg:   fieldError.Translate(m.Tran),
			}

			// get original tag name from value for set err field key.
			structNamespace := fieldError.StructNamespace()
			_, fieldName, found := strings.Cut(structNamespace, ".")
			if found {
				originalTag := getObjectTagByFieldName(value, fieldName)
				if len(originalTag) > 0 {
					errField.ErrorField = originalTag
				}
			}
			errFields = append(errFields, errField)
		}
		if len(errFields) > 0 {
			errMsg := ""
			if len(errFields) == 1 {
				errMsg = errFields[0].ErrorMsg
			}
			return errFields, errors.BadRequest(reason.RequestFormatError).WithMsg(errMsg)
		}
	}

	if v, ok := value.(Checker); ok {
		errFields, err = v.Check()
		if err == nil {
			return nil, nil
		}
		errMsg := ""
		for _, errField := range errFields {
			errField.ErrorMsg = translator.Tr(m.Lang, errField.ErrorMsg)
			errMsg = errField.ErrorMsg
		}
		return errFields, errors.BadRequest(reason.RequestFormatError).WithMsg(errMsg)
	}
	return nil, nil
}

// Checker .
type Checker interface {
	Check() (errField []*FormErrorField, err error)
}

func getObjectTagByFieldName(obj interface{}, fieldName string) (tag string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	objT := reflect.TypeOf(obj)
	objT = objT.Elem()

	structField, exists := objT.FieldByName(fieldName)
	if !exists {
		return ""
	}
	tag = structField.Tag.Get("json")
	if len(tag) == 0 {
		return structField.Tag.Get("form")
	}
	return tag
}
