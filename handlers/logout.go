package handlers

import (
	"net/http"
)

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the session
		session, _ := store.Get(r, "store") // Replace "session-name" with your actual session name

		// Clear specific session values
		session.Values["sturecmsaid"] = nil
		session.Values["username"] = nil
		session.Values["adm"] = nil
		session.Values["phone"] = nil
		session.Values["password"] = nil

		// Save the session after clearing
		session.Save(r, w)

		// Redirect to the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
