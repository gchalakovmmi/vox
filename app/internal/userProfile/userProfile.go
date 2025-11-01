package userProfile

import (
	"net/mail"
	"unicode"
	"golang.org/x/crypto/bcrypt"
)

type UserProfile struct {
	ID			int
	FirstName		string
	LastName		string
	Username		string
	Email			string
	Password		string
	RepeatedPassword	string
	HashedPassword		string
}

type PasswordStrength struct {
	Length	bool
	Lower	bool
	Upper	bool
	Number	bool
	Special	bool
}

func (u *UserProfile) verifyPassword() PasswordStrength {
	var strength PasswordStrength
	for _, c := range u.Password {
		switch {
		case unicode.IsLower(c):
			strength.Lower = true
		case unicode.IsUpper(c):
			strength.Upper = true
		case unicode.IsNumber(c):
			strength.Number = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			strength.Special = true
		case unicode.IsLetter(c) || c == ' ':
		default:
			return PasswordStrength{}
		}
	}
	strength.Length = len(u.Password) >= 8
	return strength
}

func (u *UserProfile) ValidateSignUpInfo() string {
	if u.FirstName == "" {
		return "signup_first_name"
	}

	if u.LastName == "" {
		return "signup_last_name"
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return "signup_email"
	}

	strength := u.verifyPassword()
	switch {
	case !strength.Length:
		return "signup_password_length"
	case !strength.Upper:
		return "signup_password_uppercase"
	case !strength.Lower:
		return "signup_password_lowercase"
	case !strength.Number:
		return "signup_password_number"
	case !strength.Special:
		return "signup_password_special"
	}

	if u.Password != u.RepeatedPassword {
		return "signup_repeated_password"
	}

	return ""
}

func (u *UserProfile) HashPassword() (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return "", err
	}
	u.HashedPassword = string(hash)
	return u.HashedPassword, nil
}

func (u *UserProfile) CompareHashAndPassword() error{
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(u.Password))
}
