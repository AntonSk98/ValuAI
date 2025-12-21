package common

// Language represents a supported language code in the system.
type Language string

// Supported language codes.
const (
	En Language = "en" // English
	Pl Language = "pl" // Polish
	De Language = "de" // German
)

// IsLanguageSupported checks if the given Language is supported
func IsLanguageSupported(lang Language) bool {
	switch lang {
	case En, Pl, De:
		return true
	default:
		return false
	}
}
