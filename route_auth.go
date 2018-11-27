package main

import (
	"github.com/halberdholder/bbs/data"
	"net/http"
	"strconv"
)

type LoginInfo struct {
	ThreadUuid string
	Failed     bool
	Admin      bool
}

func (loginInfo LoginInfo) Redirect(w http.ResponseWriter) {
	generateHTML(w, loginInfo, "login.layout", "login")
}

// GET /login
// Show the login page
func login(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}

	isFailed, _ := strconv.ParseBool(request.FormValue("failed"))

	loginInfo := LoginInfo{
		ThreadUuid: request.FormValue("uuid"),
		Failed:     isFailed,
		Admin:      false,
	}

	generateHTML(writer, loginInfo, "login.layout", "login")
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
		generateHTML(writer, true, "login.layout", "public.navbar", "signup")
		danger(err, "Cannot create user")
	} else {
		http.Redirect(writer, request, "/login", 302)
	}
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	threadUuid := request.PostFormValue("uuid")
	loginInfo := LoginInfo{
		ThreadUuid: threadUuid,
		Failed:     true,
	}
	user, err := data.UserByEmail(request.PostFormValue("email"))
	if err != nil {
		loginInfo.Redirect(writer)
		danger(err, "Cannot find user")
		return
	}

	if isAdministrator(user) {
		loginInfo.Redirect(writer)
		info("Administrotor cannot login normal site")
		return
	}

	if user.Password != data.Encrypt(request.PostFormValue("password")) {
		loginInfo.Redirect(writer)
		info("user", user.Email, "password error, login failed")
		return
	}

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

	switch threadUuid {
	case "new":
		http.Redirect(writer, request, "/thread/new", 302)
	case "":
		http.Redirect(writer, request, "/", 302)
	default:
		http.Redirect(writer, request, "/thread/read?id="+threadUuid, 302)
	}
	info("user", user.Email, "login success")
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
