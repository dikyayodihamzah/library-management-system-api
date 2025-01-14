package lib

import "regexp"

// IsValidEmail check is valid email in given string
func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}

	// Define the email regex pattern
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	// Match the email against the regex pattern
	return re.MatchString(email)
}
