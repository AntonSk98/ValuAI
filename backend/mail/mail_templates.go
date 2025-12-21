package mail

import (
	"errors"
	"fmt"
	"strings"
	"valuai/common"
)

// templateContent holds the title and body of a template for one particular language.
type templateContent struct {
	Title string `yaml:"title"`
	Body  string `yaml:"body"`
}

// MailTemplates represents the loaded YAML structure
type MailTemplates struct {
	Templates map[string]map[string]templateContent `yaml:"templates"`
}

// Global templates instance
var mailTemplates *MailTemplates

// ResolveTemplates loads mail templates from the specified YAML file path
func ResolveTemplates(path string) {
	var templates MailTemplates
	err := common.LoadProperties(path, &templates)
	if err != nil {
		panic(fmt.Sprintf("failed to load mail templates: %v", err))
	}
	mailTemplates = &templates
}

// GetTemplate fetches a template content by its email template name and language.
// If the template name or language is unsupported, an error is returned.
func GetTemplate(name TemplateName, lang common.Language) (*templateContent, error) {
	if mailTemplates == nil {
		return nil, fmt.Errorf("mail templates not loaded")
	}

	templateData, ok := mailTemplates.Templates[string(name)]
	if !ok {
		return nil, fmt.Errorf("unsupported template: %s", name)
	}

	content, ok := templateData[string(lang)]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	return &content, nil
}

// ResolveTemplate fills the %s placeholders in the email template body.
// The placeholders are filled in the order of the provided arguments.
// Returns an error if there are unresolved placeholders remaining.
func (t *templateContent) ResolveTemplateContent(args ...any) (string, error) {
	result := fmt.Sprintf(t.Body, args...)
	if strings.Contains(result, "%s") {
		return "", errors.New("unresolved placeholders remain in template")
	}
	return result, nil
}
