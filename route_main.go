package main

import (
	"fmt"
	"github.com/halberdholder/bbs/data"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// GET /err?msg=
// shows the error message page
func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}

func index(writer http.ResponseWriter, request *http.Request) {
	threads, err := data.Threads()
	if err != nil {
		error_message(writer, request, "Cannot get threads")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, threads, "layout", "public.navbar", "index")
			info("anonymous", request.Host, "visited")
		} else {
			generateHTML(writer, threads, "layout", "private.navbar", "index")
		}
	}
}

func ckeditorUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseMultipartForm(4 << 20) // 4M
		if err != nil {
			danger(err)
			http.Error(w, err.Error(), 500)
			return
		}
		mpartFile, mpartHeader, err := r.FormFile("upload")
		if err != nil {
			danger(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer mpartFile.Close()
		uri, err := saveFile(mpartHeader, mpartFile)
		if err != nil {
			danger(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		CKEdFunc := r.FormValue("CKEditorFuncNum")
		fmt.Fprintln(w, "<script>window.parent.CKEDITOR.tools.callFunction("+CKEdFunc+", \""+uri+"\");</script>")

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		danger(err)
		http.Error(w, err.Error(), 405)
	}
}

func saveFile(fh *multipart.FileHeader, f multipart.File) (string, error) {
	fileExt := filepath.Ext(fh.Filename)
	newName := fmt.Sprint(time.Now().Unix()) + fileExt //unique file name based on timestamp. You can keep original name and ignore duplicates
	uri := "/static/uploads/images/" + newName
	fullName := filepath.Join("./public/uploads/images/", newName)

	file, err := os.OpenFile(fullName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, f)
	if err != nil {
		return "", err
	}
	return uri, nil
}
