package utils

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword Basic password policy, At least 8 characters, contains at least one letter and one number
func ValidatePassword(password string) bool {
	passwordRegex := regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`)
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)

	return passwordRegex.MatchString(password) && hasLetter && hasNumber
}
