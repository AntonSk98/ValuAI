package mail

import (
	"errors"
	"fmt"
	"strings"
	"valuai/common"
)

// templateContent holds the title and body of a template for one paricular language.
type templateContent struct {
	title string
	body  string
}

// template represents an email template localized for each of the supported languages.
type template struct {
	languages map[common.Language]templateContent
}

// emailTemplates holds all supported email templates.
type emailTemplates struct {
	templates map[TemplateNames]template
}

// internal constant with all email templates.
var templates = emailTemplates{
	templates: map[TemplateNames]template{
		OtpEmail: {
			languages: map[common.Language]templateContent{
				common.En: {
					title: "🔐 ValuAI: Your login code",
					body:  "Hello!\n\nYour one-time login code is: *%s*\n\nPlease enter this code in the app to continue. The code is valid for only a few minutes.\n\nIf you didn’t request this code, you can safely ignore this message.",
				},
				common.Pl: {
					title: "🔐 ValuAI: Twój kod logowania",
					body:  "Cześć!\n\nTwój jednorazowy kod logowania to: *%s*\n\nWpisz ten kod w aplikacji, aby kontynuować. Kod jest ważny tylko przez kilka minut.\n\nJeśli nie prosiłeś o ten kod, zignoruj tę wiadomość.",
				},
				common.De: {
					title: "🔐 ValuAI: Dein Anmeldecode",
					body:  "Hallo!\n\nDein einmaliger Anmeldecode lautet: *%s*\n\nBitte gib diesen Code in der App ein, um fortzufahren. Der Code ist nur wenige Minuten gültig.\n\nWenn du diesen Code nicht angefordert hast, kannst du diese Nachricht ignorieren.",
				},
			},
		},
	},
}

// GetTemplate fetches a template content by its email template name and language.
// If the template name or language is unsupported, an error is returned.
func GetTemplate(name TemplateNames, lang common.Language) (*templateContent, error) {
	tmpl, ok := templates.templates[name]
	if !ok {
		return nil, fmt.Errorf("unsupported template: %s", name)
	}

	content, ok := tmpl.languages[lang]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	return &content, nil
}

// ResolveTemplate fills the %s placeholders in the email template body.
// The placeholders are filled in the order of the provided arguments.
// Returns an error if there are unresolved placeholders remaining.
func (t *templateContent) ResolveTemplateContent(args ...any) (string, error) {
	result := fmt.Sprintf(t.body, args...)
	if strings.Contains(result, "%s") {
		return "", errors.New("unresolved placeholders remain in template")
	}
	return result, nil
}

// Optional: also expose Title getter
func (t *templateContent) Title() string {
	return t.title
}
