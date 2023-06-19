package main

import (
	"net/http"
	"fmt"
	
	"github.com/gorilla/sessions"
    "github.com/julienschmidt/httprouter"
)


var store = sessions.NewCookieStore([]byte("secret_key"))


func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		fmt.Println("path", r.URL.Path)
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "register.html", nil)
	} else {
			fmt.Println("path", r.URL.Path)
			r.ParseForm()
			username := r.Form.Get("username")
			password := r.Form.Get("password")
			email := r.Form.Get("email")
			fmt.Println("RDX says: ", rdxHset("userlogin",username,password))
			fmt.Println("method:", r.Method)
			fmt.Println(username, password, email)

			rdxuser:= rdxHget("userlogin",username)
			if username == rdxuser {
				println("User already exists")
			}
			if password !="" && username != "" && email != "" {
				fmt.Println("Success")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			}
		}
	}



func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		fmt.Println("path", r.URL.Path)
		if isLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)			
		}
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "login.html", nil)
	} else {
			fmt.Println("path", r.URL.Path)
			r.ParseForm()
			username := r.Form.Get("username")
			password := r.Form.Get("password")
			fmt.Println("method:", r.Method)
			fmt.Println(username, password)
			session, _ := store.Get(r, "session.id")
			rdxpass:= rdxHget("userlogin",username)
			fmt.Println("RDX pass: ", rdxpass)
			if password ==rdxpass && password != "" {
				session.Values["authenticated"] = true 
				session.Save(r, w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		}
	}





func logoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}
	session, _ := store.Get(r, "session.id")
	session.Values["authenticated"] = false 
	session.Save(r, w)
//	w.Write([]byte("Logout Successful"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}






func Dash(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Welcome!"))
}

func isLoggedIn(r *http.Request) bool {
	session, _ := store.Get(r, "session.id")
	authenticated := session.Values["authenticated"]
	if authenticated != nil && authenticated != false {
		return true
	}
	return false
}

