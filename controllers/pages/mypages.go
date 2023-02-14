package pages

import (
	"html/template"
	"net/http"
	"projek/entities"

	"github.com/gorilla/sessions"
)

func AddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("./views/projectAdd.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	if session.Values["IsLogin"] != true {
		entities.Data["IsLogin"] = false
		entities.Data["Id"] = session.Values["Id"]
} else {
		entities.Data["IsLogin"] = session.Values["IsLogin"].(bool)
		entities.Data["UserName"] = session.Values["Name"].(string)
		entities.Data["Id"] = session.Values["Id"].(int)
}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, entities.Data)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("./views/contact-me.html")
	// error handling
	if err != nil {
		panic(err)
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		entities.Data["IsLogin"] = false
		entities.Data["Id"] = session.Values["Id"]
	} else {
		entities.Data["IsLogin"] = session.Values["IsLogin"].(bool)
		entities.Data["UserName"] = session.Values["Name"].(string)
		entities.Data["Id"] = session.Values["Id"].(int)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, entities.Data)
}

func FormLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("./views/pageLogin.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var data = map[string]interface{}{
		"title":   "Login | Marcel",
		"isLogin": true,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func FormRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("./views/pageRegister.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var data = map[string]interface{}{
		"title":   "Register | Marcel",
		"isLogin": true,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}
