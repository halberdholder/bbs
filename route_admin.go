package main

import (
	"github.com/halberdholder/bbs/data"
	"net/http"
	"strconv"
)

type Permission int

const (
	PermAddThreadClass = 1 << 0
	PermModThreadClass = 1 << 1
	PermDelThreadClass = 1 << 2
	PermDelThread      = 1 << 3
)

type AdminPageInfo struct {
	PageInfo
	Perm Permission
}

func adminErr(writer http.ResponseWriter, request *http.Request) {
	if _, err := session(request, "admin_cookie"); err != nil {
		redirectLogin(writer)
		return
	}
	vals := request.URL.Query()
	generateHTML(writer, vals.Get("msg"), "layout", "admin.navbar", "error")
}

func adminLogin(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}

	isFailed, _ := strconv.ParseBool(request.FormValue("failed"))

	loginInfo := LoginInfo{
		Failed: isFailed,
		Admin:  true,
	}

	generateHTML(writer, loginInfo, "login.layout", "login")
}

func adminAuthenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	user, err := data.UserByEmail(request.PostFormValue("email"))
	if err != nil || !isAdministrator(user) {
		redirectLogin(writer)
		return
	}

	if user.Password != data.Encrypt(request.PostFormValue("password")) {
		redirectLogin(writer)
		return
	}

	session, err := user.CreateSession()
	if err != nil {
		danger(err, "Cannot create session")
		redirectLogin(writer)
		return
	}

	cookie := http.Cookie{
		Name:  "admin_cookie",
		Value: session.Uuid,
	}
	http.SetCookie(writer, &cookie)

	http.Redirect(writer, request, "/index", 302)

	info("admin", user.Email, "login success")
}

func adminIndex(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(request, "admin_cookie")
	if err != nil {
		redirectLogin(writer)
		return
	}
	user, err := sess.User()
	if err != nil {
		redirectLogin(writer)
		return
	}

	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}
	currentPage, err := strconv.ParseInt(request.FormValue("page"), 10, 64)
	if currentPage == 0 || err != nil {
		currentPage = 1
	}

	adminPageInfo := AdminPageInfo{
		PageInfo: PageInfo{
			CurrentPage: currentPage,
		},
	}
	total, err := data.TotalThreads()
	if err != nil {
		error_message(writer, request, "Cannot get total count of threads")
		return
	}
	adminPageInfo.TotalThreads = total

	threads, err := data.ThreadsByPage(currentPage, config.PageSize)
	if err != nil {
		error_message(writer, request, "Cannot get threads")
		return
	}
	adminPageInfo.Threads = threads

	adminPageInfo.Pagination()
	adminPageInfo.Perm = Permission(user.Permission)
	generateHTML(writer, adminPageInfo, "layout", "admin.navbar", "admin.index", "thread.page")
}

func adminClass(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(request, "admin_cookie")
	if err != nil {
		redirectLogin(writer)
		return
	}
	user, err := sess.User()
	if err != nil {
		redirectLogin(writer)
		return
	}

	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}
	currentPage, err := strconv.ParseInt(request.FormValue("page"), 10, 64)
	if currentPage == 0 || err != nil {
		currentPage = 1
	}
	adminPageInfo := AdminPageInfo{
		PageInfo: PageInfo{
			CurrentPage: currentPage,
		},
	}

	threadClassId, err := strconv.ParseInt(request.FormValue("id"), 10, 64)
	if threadClassId <= 0 || err != nil {
		http.Redirect(writer, request, "/index", 302)
		return
	}
	adminPageInfo.ThreadClass = threadClassId

	total, err := data.TotalThreadsOfClass(adminPageInfo.ThreadClass)
	if err != nil {
		error_message(writer, request, "Cannot get total count of threads")
		return
	}
	adminPageInfo.TotalThreads = total

	threads, err := data.ThreadsByClassAndPage(adminPageInfo.ThreadClass, currentPage, config.PageSize)
	if err != nil {
		error_message(writer, request, "Cannot get threads")
		return
	}
	adminPageInfo.Threads = threads

	adminPageInfo.Pagination()
	adminPageInfo.Perm = Permission(user.Permission)
	generateHTML(writer, adminPageInfo, "layout", "admin.navbar", "admin.index", "threadclass.page")
}

func addThreadClass(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(request, "admin_cookie")
	if err != nil {
		redirectLogin(writer)
		return
	}
	user, _ := sess.User()
	if !Permission(user.Permission).CanAddThreadClass() {
		redirectLogin(writer)
		return
	}

	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}
	err = data.AddThreadClass(request.FormValue("name"))
	if err != nil {
		error_message(writer, request, err.Error())
	} else {
		http.Redirect(writer, request, "/index", 302)
	}
}

func delThreadClass(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(request, "admin_cookie")
	if err != nil {
		redirectLogin(writer)
		return
	}
	user, _ := sess.User()
	if !Permission(user.Permission).CanDelThreadClass() {
		redirectLogin(writer)
		return
	}

	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}
	err = data.DeleteThreadClassById(request.FormValue("id"))
	if err != nil {
		error_message(writer, request, err.Error())
	} else {
		http.Redirect(writer, request, "/index", 302)
	}
}

func modThreadClass(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(request, "admin_cookie")
	if err != nil {
		redirectLogin(writer)
		return
	}
	user, _ := sess.User()
	if !Permission(user.Permission).CanModThreadClass() {
		redirectLogin(writer)
		return
	}

	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}
	err = data.ModifyThreadClass(request.FormValue("id"), request.FormValue("class"))
	if err != nil {
		error_message(writer, request, err.Error())
	} else {
		http.Redirect(writer, request, "/index", 302)
	}
}

func delThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(request, "admin_cookie")
	if err != nil {
		redirectLogin(writer)
		return
	}
	user, _ := sess.User()
	if !Permission(user.Permission).CanDelThread() {
		redirectLogin(writer)
		return
	}

	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}
	err = data.DeleteThread(request.FormValue("id"))
	if err != nil {
		error_message(writer, request, err.Error())
	} else {
		http.Redirect(writer, request, "/index", 302)
	}
}

func adminLogout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("admin_cookie")
	if err != http.ErrNoCookie {
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	} else {
		warning(err, "Failed to get cookie")
	}

	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}
	if request.FormValue("index") == "true" {
		http.Redirect(writer, request, "http://"+config.PubAddr, 302)
	} else {
		http.Redirect(writer, request, "/", 302)
	}
}

func (p Permission) CanAddThreadClass() bool {
	return 0 != (p & PermAddThreadClass)
}

func (p Permission) CanModThreadClass() bool {
	return 0 != (p & PermModThreadClass)
}

func (p Permission) CanDelThreadClass() bool {
	return 0 != (p & PermDelThreadClass)
}

func (p Permission) CanDelThread() bool {
	return 0 != (p & PermDelThread)
}

func isAdministrator(user data.User) bool {
	p := Permission(user.Permission)
	return p.CanAddThreadClass() || p.CanModThreadClass() || p.CanDelThreadClass() || p.CanDelThread()
}

func redirectLogin(w http.ResponseWriter) {
	loginInfo := LoginInfo{
		Failed: true,
		Admin:  true,
	}
	loginInfo.Redirect(w)
}

func PubAddress() string {
	return config.PubAddr
}
