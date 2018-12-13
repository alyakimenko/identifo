package html

import (
	"net/http"

	"github.com/madappgang/identifo/model"
)

//ResetPassword handles password reset form submition
func (ar *Router) ResetPassword() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		password := r.FormValue("password")
		if err := model.StrongPswd(password); err != nil {
			SetFlash(w, FlashErrorMessageKey, err.Error())
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		tokenString := r.Context().Value(model.TokenRawContextKey).(string)
		token, _ := ar.TokenService.Parse(tokenString)

		err := ar.UserStorage.ResetPassword(token.UserID(), password)
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		http.Redirect(w, r, "./reset/success", http.StatusMovedPermanently)
	}

}
