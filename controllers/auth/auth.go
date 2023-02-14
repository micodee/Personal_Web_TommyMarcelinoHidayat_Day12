package auth

import (
	"context"
	"log"
	"net/http"
	"projek/config"
	"projek/entities"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("pass")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	postUSER := "INSERT INTO public.tb_users(name, email, password) VALUES ($1, $2, $3);"
	_ , err = config.ConnDB.Exec(context.Background(), postUSER, name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}



func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password :=r.PostForm.Get("pass")
	user := entities.Users{}
	selectUSER := "SELECT id, name, email, password FROM tb_users WHERE email=$1"
	rows := config.ConnDB.QueryRow(context.Background(), selectUSER, email)
	err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 10800

	// flash login
	session.AddFlash("Login Success", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-chace, no-store, must-revalidate")

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Options.MaxAge = -1

	session.AddFlash("You have successfully logged out", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}