package main

import (
	"github.com/halberdholder/bbs/data"
	"net/http"
	"strconv"
)

type LoginInfo struct {
	ThreadUuid string
	Failed bool
}

// GET /login
// Show the login page
func login(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}

	failed, _ := strconv.ParseBool(request.FormValue("failed"))

	loginInfo := LoginInfo {
		ThreadUuid: request.FormValue("uuid"),
		Failed: failed,
	}

	generateHTML(writer, loginInfo,"login.layout", "public.navbar", "login")
}

// GET /signup
// Show the signup page
func signup(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup
// Create the user account
func signupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	user := data.User{
		Name:     request.PostFormValue("name"),
		Email:    request.PostFormValue("email"),
		Password: request.PostFormValue("password"),
	}
	info("user signup", user)
	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}
	http.Redirect(writer, request, "/login", 302)
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	user, err := data.UserByEmail(request.PostFormValue("email"))
	if err != nil {
		danger(err, "Cannot find user")
	}
	if user.Password == data.Encrypt(request.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)

		threadUuid := request.PostFormValue("uuid")
		if threadUuid != ""{
			http.Redirect(writer, request, "/thread/read?id="+"threadUuid", 302)
		} else {
			http.Redirect(writer, request, "/", 302)
		}
		info("user", user.Email, "login success")
	} else {
		http.Redirect(writer, request, "/login?failed=true", 302)
		info("user", user.Email, "login failed")
	}

}

// GET /logout
// Logs the user out
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	} else {
		warning(err, "Failed to get cookie")
	}
	http.Redirect(writer, request, "/", 302)
}
