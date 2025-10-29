package handlers

import (
	"log/slog"
	"net/http"
	"net/mail"
	"net/url"
	"unicode"
	"golang.org/x/crypto/bcrypt"
)

type UserProfile struct {
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

func (u *UserProfile) validateSignUpInfo() string {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return "email"
	}

	strength := u.verifyPassword()
	switch {
	case !strength.Length:
		return "password_length"
	case !strength.Upper:
		return "password_uppercase"
	case !strength.Lower:
		return "password_lowercase"
	case !strength.Number:
		return "password_number"
	case !strength.Special:
		return "password_special"
	}

	if u.Password != u.RepeatedPassword {
		return "repeatedPassword"
	}

	return ""
}

func (u *UserProfile) HashPassword() error{
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hash)
	return nil
}

func (u *UserProfile) ComparePasswordAndHash() error{
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(u.HashedPassword))
}

func ProcessSignup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userProfile := UserProfile{
			Email:			r.PostFormValue("email"),
			Password:		r.PostFormValue("password"),
			RepeatedPassword:	r.PostFormValue("password-repeat"),
		}

		queryParams := url.Values{}
		queryParams.Add("Email", url.QueryEscape(userProfile.Email))

		if errorMessage := userProfile.validateSignUpInfo(); errorMessage != "" {
			queryParams.Add("error", errorMessage)
			http.Redirect(w, r, "/signup?"+queryParams.Encode(), http.StatusSeeOther)
			return
		}

		slog.Debug("SignUp info Correct! internal/handlers/process-signup.go", "Email", userProfile.Email, "Password", userProfile.Password, "Repeated Password", userProfile.RepeatedPassword)

		if userProfile.HashPassword() != nil {
			queryParams.Add("error", "hashing")
			http.Redirect(w, r, "/signup?"+queryParams.Encode(), http.StatusSeeOther)
			return
		}
		slog.Debug("Hashing Successful", "Hashed Password", userProfile.HashedPassword)
	})
}
