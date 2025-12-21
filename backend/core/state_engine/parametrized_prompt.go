package core

import "strings"

// ParametrizedPrompt represents a prompt with placeholders.
type ParametrizedPrompt struct {
	LocalizedPrompt string
}

// AsParametrizedPrompt creates a new ParametrizedPrompt instance.
func AsParametrizedPrompt(localizedPrompt string) *ParametrizedPrompt {
	return &ParametrizedPrompt{
		LocalizedPrompt: localizedPrompt,
	}
}

// Resolve replaces all occurrences of "@key@" with the given value.
// Returns the same object for chaining.
func (p *ParametrizedPrompt) Resolve(key string, value string) *ParametrizedPrompt {
	p.LocalizedPrompt = strings.ReplaceAll(p.LocalizedPrompt, "@"+key+"@", value)
	return p
}

// Render returns the final prompt string after all replacements.
func (p *ParametrizedPrompt) Render() string {
	return p.LocalizedPrompt
}
