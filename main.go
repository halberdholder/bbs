package main

import (
	"net/http"
	"time"
)

func httpServer() {
	//func httpserver(wg *sync.WaitGroup) {
	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//
	// all route patterns matched here
	// route handler functions defined in other files
	//

	// index
	mux.HandleFunc("/", index)
	mux.HandleFunc("/index", index)
	// error
	mux.HandleFunc("/err", err)

	mux.HandleFunc("/ckeditorUpload", ckeditorUpload)

	// defined in route_auth.go
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/authenticate", authenticate)

	// defined in route_thread.go
	mux.HandleFunc("/thread/new", newThread)
	mux.HandleFunc("/thread/create", createThread)
	mux.HandleFunc("/thread/post", postThread)
	mux.HandleFunc("/thread/read", readThread)

	mux.HandleFunc("/class", threadClass)

	// starting up the server
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()

	//wg.Done()
}

func httpsServer() {
	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//
	// all route patterns matched here
	// route handler functions defined in other files
	//

	// index and login
	mux.HandleFunc("/", adminLogin)
	// error
	mux.HandleFunc("/err", adminErr)

	// defined in route_auth.go
	mux.HandleFunc("/authenticate", adminAuthenticate)
	mux.HandleFunc("/index", adminIndex)
	mux.HandleFunc("/class", adminClass)
	mux.HandleFunc("/logout", adminLogout)

	mux.HandleFunc("/addthreadclass", addThreadClass)
	mux.HandleFunc("/delthreadclass", delThreadClass)
	mux.HandleFunc("/modthreadclass", modThreadClass)
	mux.HandleFunc("/delthread", delThread)
	/*

		mux.HandleFunc("/signup", signup)
		mux.HandleFunc("/signup_account", signupAccount)

		// defined in route_thread.go
		mux.HandleFunc("/thread/new", newThread)
		mux.HandleFunc("/thread/create", createThread)
		mux.HandleFunc("/thread/post", postThread)
		mux.HandleFunc("/thread/read", readThread)

		mux.HandleFunc("/class", threadClass)
	*/

	// starting up the server
	server := &http.Server{
		Addr:           config.TLSAddress,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServeTLS(config.TLSCert, config.TLSKey)
}

func main() {
	p("ChitChat", version(), "started at", config.Address)

	go httpServer()
	go httpServer()

	select {}
}
