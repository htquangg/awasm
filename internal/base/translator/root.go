// https://github.com/apache/incubator-answer/blob/main/internal/base/translator/provider.go
package translator

import (
	"fmt"
	"os"
	"path/filepath"

	myTran "github.com/segmentfault/pacman/contrib/i18n"
	"github.com/segmentfault/pacman/i18n"
	"gopkg.in/yaml.v3"

	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/pkg/logger"
)

var GlobalTrans i18n.Translator

// LangOption language option
type LangOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
	// Translation completion percentage
	Progress int `json:"progress"`
}

// DefaultLangOption default language option. If user config the language is default, the language option is admin choose.
const DefaultLangOption = "Default"

// LanguageOptions language
var LanguageOptions []*LangOption

// NewTranslator new a translator
func NewTranslator(c *config.I18n) (tr i18n.Translator, err error) {
	entries, err := os.ReadDir(c.BundleDir)
	if err != nil {
		return nil, err
	}

	// read the Bundle resources file from entries
	for _, file := range entries {
		// ignore directory
		if file.IsDir() {
			continue
		}
		// ignore non-YAML file
		if filepath.Ext(file.Name()) != ".yaml" && file.Name() != "i18n.yaml" {
			continue
		}
		logger.Debugf("try to read file: %s", file.Name())
		buf, err := os.ReadFile(filepath.Join(c.BundleDir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("read file failed: %s %s", file.Name(), err)
		}

		// parse the backend translation
		originalTr := struct {
			Backend map[string]map[string]interface{} `yaml:"backend"`
			UI      map[string]interface{}            `yaml:"ui"`
		}{}
		if err = yaml.Unmarshal(buf, &originalTr); err != nil {
			return nil, err
		}
		translation := make(map[string]interface{}, 0)
		for k, v := range originalTr.Backend {
			translation[k] = v
		}
		translation["backend"] = originalTr.Backend
		translation["ui"] = originalTr.UI

		content, err := yaml.Marshal(translation)
		if err != nil {
			logger.Debugf("marshal translation content failed: %s %s", file.Name(), err)
			continue
		}

		// add translator use backend translation
		if err = myTran.AddTranslator(content, file.Name()); err != nil {
			logger.Debugf("add translator failed: %s %s", file.Name(), err)
			continue
		}
	}
	GlobalTrans = myTran.GlobalTrans

	i18nFile, err := os.ReadFile(filepath.Join(c.BundleDir, "i18n.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read i18n file failed: %s", err)
	}

	s := struct {
		LangOption []*LangOption `yaml:"language_options"`
	}{}
	err = yaml.Unmarshal(i18nFile, &s)
	if err != nil {
		return nil, fmt.Errorf("i18n file parsing failed: %s", err)
	}
	LanguageOptions = s.LangOption
	for _, option := range LanguageOptions {
		option.Label = fmt.Sprintf("%s (%d%%)", option.Label, option.Progress)
	}
	return GlobalTrans, err
}

// CheckLanguageIsValid check user input language is valid
func CheckLanguageIsValid(lang string) bool {
	if lang == DefaultLangOption {
		return true
	}
	for _, option := range LanguageOptions {
		if option.Value == lang {
			return true
		}
	}
	return false
}

// Tr use language to translate data. If this language translation is not available, return default english translation.
func Tr(lang i18n.Language, data string) string {
	if GlobalTrans == nil {
		return data
	}
	translation := GlobalTrans.Tr(lang, data)
	if translation == data {
		return GlobalTrans.Tr(i18n.DefaultLanguage, data)
	}
	return translation
}

// TrWithData translate key with template data, it will replace the template data {{ .PlaceHolder }} in the translation.
func TrWithData(lang i18n.Language, key string, templateData any) string {
	if GlobalTrans == nil {
		return key
	}
	translation := GlobalTrans.TrWithData(lang, key, templateData)
	if translation == key {
		return GlobalTrans.TrWithData(i18n.DefaultLanguage, key, templateData)
	}
	return translation
}
