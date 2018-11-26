package main

import (
	"net/http"
	"strconv"
)

func adminlogin(writer http.ResponseWriter, request *http.Request) {
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
