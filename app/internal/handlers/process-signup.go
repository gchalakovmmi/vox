package handlers

import (
	up "vox/internal/userProfile"
	"log/slog"
	"net/http"
	"net/url"
)

func ProcessSignup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userProfile := up.UserProfile{
			FirstName:		r.PostFormValue("first_name"),
			LastName:		r.PostFormValue("last_name"),
			Email:			r.PostFormValue("email"),
			Password:		r.PostFormValue("password"),
			RepeatedPassword:	r.PostFormValue("password-repeat"),
		}

		queryParams := url.Values{}
		queryParams.Add("first_name", url.QueryEscape(userProfile.FirstName))
		queryParams.Add("last_name", url.QueryEscape(userProfile.LastName))
		queryParams.Add("email", url.QueryEscape(userProfile.Email))

		if errorMessage := userProfile.ValidateSignUpInfo(); errorMessage != "" {
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
