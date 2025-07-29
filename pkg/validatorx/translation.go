package validatorx

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/epkgs/i18n"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
	"golang.org/x/text/language"
)

const (
	VALIDATE_DateAfterToday = "dateAfterToday"
)

var g_supported = []language.Tag{
	language.Make("en"),
	language.Make("zh"),
}
var g_matcher = language.NewMatcher(g_supported)

var g_trans = map[string]ut.Translator{}

var g_uni_translator *ut.UniversalTranslator

func init() {

	zhT := zh.New() // 中文翻译器
	enT := en.New() // 英文翻译器

	// 第一个参数是备用（fallback）的语言环境
	// 后面的参数是应该支持的语言环境（支持多个）
	g_uni_translator = ut.New(enT, zhT, enT)
}

func Translator(locale ...string) ut.Translator {

	tags := []language.Tag{}
	for _, l := range locale {
		tags = append(tags, language.Make(l))
	}

	_, i, _ := g_matcher.Match(tags...)
	lang := g_supported[i]
	lc := lang.String()

	if _, ok := g_trans[lc]; !ok {
		g_trans[lc], _ = NewTranslator(lc)
	}

	return g_trans[lc]
}

func TranslatorDetect(ctx context.Context) ut.Translator {
	langs := i18n.GetAcceptLanguages(ctx)
	return Translator(langs...)
}

// ! 无法翻译错误信息内的 attribute 为中文
func NewTranslator(locale string) (trans ut.Translator, err error) {

	// 修改gin框架中的Validator引擎属性，实现自定制
	if engine, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json tag的自定义方法
		engine.RegisterTagNameFunc(func(fld reflect.StructField) string {

			// TODO: 修改描述内的字段名称
			jsonTag := fld.Tag.Get("json")

			if jsonTag == "" {
				return ""
			}

			name := strings.SplitN(jsonTag, ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		trans, _ = g_uni_translator.GetTranslator(locale)

		// 注册翻译器
		switch trans.Locale() {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(engine, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(engine, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(engine, trans)
		}

		if err != nil {
			return
		}

		// 在校验器注册自定义的校验方法
		if err = engine.RegisterValidation(VALIDATE_DateAfterToday, validateDateAfterToday); err != nil {
			return
		}

		// 注册自定义的校验字段
		if err = engine.RegisterTranslation(
			VALIDATE_DateAfterToday,
			trans,
			registerTranslator(VALIDATE_DateAfterToday, "{0}必须要晚于当前日期"),
			translate,
		); err != nil {
			return
		}
	}
	err = errors.New("failed to create new translator")
	return

}

// validateDateAfterToday 日期校验器,校验传入time.Date是否晚于当前时间
func validateDateAfterToday(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	if date.Before(time.Now()) {
		return false
	}
	return true
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}
