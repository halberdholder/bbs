package main

import (
	"net/http"
	"strconv"
	"github.com/halberdholder/bbs/data"
)

func threadClass(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		danger(err, "Cannot prase form")
	}

	currentPage, err := strconv.ParseInt(request.FormValue("page"), 10, 64)
	if currentPage == 0 || err != nil {
		currentPage = 1
	}
	pageInfo := PageInfo {
		CurrentPage: currentPage,
	}

	threadClassId, err := strconv.ParseInt(request.FormValue("id"), 10, 64)
	if threadClassId <= 0 || err != nil {
		http.Redirect(writer, request, "/index", 302)
		return
	}
	pageInfo.ThreadClass = threadClassId

	total, err := data.TotalThreadsOfClass(pageInfo.ThreadClass)
	if  err != nil {
		error_message(writer, request, "Cannot get total count of threads")
		return
	}
	pageInfo.TotalThreads = total

	threads, err := data.ThreadsByClassAndPage(pageInfo.ThreadClass, currentPage, config.PageSize)
	if err != nil {
		error_message(writer, request, "Cannot get threads")
		return
	}
	pageInfo.Threads = threads

	pageInfo.Pagination()

	_, err = session(request, "_cookie")
	if err != nil {
		generateHTML(writer, pageInfo, "layout", "public.navbar", "index", "threadclass.page")
		info("anonymous", request.Host, "visited")
	} else {
		generateHTML(writer, pageInfo, "layout", "private.navbar", "index", "threadclass.page")
	}
}