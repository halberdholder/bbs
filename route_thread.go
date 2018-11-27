package main

import (
	"fmt"
	"github.com/halberdholder/bbs/data"
	"net/http"
	"strconv"
)

// GET /threads/new
// Show the new thread form page
func newThread(writer http.ResponseWriter, request *http.Request) {
	_, err := session(request, "_cookie")
	if err != nil {
		loginInfo := LoginInfo{
			ThreadUuid: "new",
		}
		generateHTML(writer, &loginInfo, "login.layout", "public.navbar", "login")
		//http.Redirect(writer, request, "/login", 302)
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "new.thread")
	}
}

// POST /signup
// Create the user account
func createThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(request, "_cookie")
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		topic := request.PostFormValue("topic")
		body := request.PostFormValue("body")
		classId, _ := strconv.ParseInt(request.PostFormValue("threadclass"), 10, 64)
		if topic != "" {
			if _, err := user.CreateThread(topic, body, classId); err != nil {
				danger(err, "Cannot create thread")
			}
		}
		http.Redirect(writer, request, "/", 302)
	}
}

// GET /thread/read
// Show the details of the thread, including the posts and the form to write a post
func readThread(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	thread, err := data.ThreadByUUID(uuid)
	if err != nil {
		error_message(writer, request, "Cannot read thread")
	} else {
		_, err := session(request, "_cookie")
		if err != nil {
			generateHTML(writer, &thread, "layout", "public.navbar", "public.thread")
		} else {
			generateHTML(writer, &thread, "layout", "private.navbar", "private.thread")
		}
	}
}

// POST /thread/post
// Create the post
func postThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(request, "_cookie")
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		uuid := request.PostFormValue("uuid")
		body := request.PostFormValue("body")
		if body != "" {
			thread, err := data.ThreadByUUID(uuid)
			if err != nil {
				error_message(writer, request, "Cannot read thread")
			}
			if _, err := user.CreatePost(thread, body); err != nil {
				danger(err, "Cannot create post")
			}
		}
		url := fmt.Sprint("/thread/read?id=", uuid)
		http.Redirect(writer, request, url, 302)
	}
}
