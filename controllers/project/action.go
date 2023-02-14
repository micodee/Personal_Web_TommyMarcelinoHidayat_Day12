package project

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"projek/config"
	"projek/entities"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")
	var tmpl, err = template.ParseFiles("./views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// session get SESSION_ID
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// flash
	fm := session.Flashes("message")
	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		} 
	}
	entities.Data["FlashData"] = strings.Join(flashes, "")
	
	// conditional is login
	var readDT string
	if session.Values["IsLogin"] != true {
		entities.Data["IsLogin"] = false
		readDT = "SELECT id, project_name, start_date, end_date, description, image, technologies, user_id FROM public.tb_projects ORDER BY id DESC;"
		} else {
			entities.Data["IsLogin"] = session.Values["IsLogin"].(bool)
			entities.Data["UserName"] = session.Values["Name"].(string)
		readDT = "SELECT tb_projects.id, project_name, start_date, end_date, description, image, technologies, tb_users.name as user FROM tb_projects LEFT JOIN tb_users ON tb_projects.user_id = tb_users.id WHERE tb_users.name='" + entities.Data["UserName"].(string) + "' ORDER BY id DESC"
	}
	rows, _ := config.ConnDB.Query(context.Background(), readDT)
					

	var result []entities.Project
 for rows.Next() {
	var each = entities.Project{}
 var err = rows.Scan(&each.Id, &each.Title, &each.Sdate, &each.Edate, &each.Content, &each.Image, &each.Technologies, &each.User)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// duration
	ConvertSdate, _ := time.Parse("02 January 2006", each.Sdate)
	ConvertEdate, _ := time.Parse("02 January 2006", each.Edate)
	duration := ConvertEdate.Sub(ConvertSdate)
	var distance string
	if duration.Hours()/24 < 7 {
		distance = strconv.FormatFloat(duration.Hours()/24, 'f', 0, 64) + " Days"
	} else if duration.Hours()/24/7 < 4 {
		distance = strconv.FormatFloat(duration.Hours()/24/7, 'f', 0, 64) + " Weeks"
	} else if duration.Hours()/24/30 < 12 {
		distance = strconv.FormatFloat(duration.Hours()/24/30, 'f', 0, 64) + " Months"
	} else {
		distance = strconv.FormatFloat(duration.Hours()/24/30/12, 'f', 0, 64) + " Years"
	}
  each.Duration = distance

		result = append(result, each)
	}

		
		resp := map[string]interface{}{
			"Data" : entities.Data,
			"Projects" : result,
		}
		
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func Add(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil{
		log.Fatal(err)
	}

	title := r.PostForm.Get("pname")
	content := r.PostForm.Get(("desc"))
	SD := r.PostForm.Get("sdate")
	ED := r.PostForm.Get("edate")
	tech := r.Form["check"]

	image := r.Context().Value("dataFile")
	img := image.(string)

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	userID := session.Values["Id"].(int)


	// convert date ke string
	sDate, _ := time.Parse("2006-01-02", SD)
	sDateFormat := sDate.Format("02 January 2006")
	eDate, _ := time.Parse("2006-01-02", ED)
	eDateFormat := eDate.Format("02 January 2006")

	addID := "INSERT INTO tb_projects(project_name, start_date, end_date, description, image, technologies, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = config.ConnDB.Exec(context.Background(), addID, title, sDateFormat, eDateFormat, content, img, tech, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}


func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("./views/projectUpdate.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	selectID := "SELECT id, project_name, start_date, end_date, description, technologies FROM tb_projects WHERE id=$1"
	rows := config.ConnDB.QueryRow(context.Background(), selectID, id)
	var getID entities.Project
	err = rows.Scan(&getID.Id, &getID.Title, &getID.Sdate, &getID.Edate, &getID.Content, &getID.Technologies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// convert string ke date
	ConvertSdate, _ := time.Parse("02 January 2006", getID.Sdate)
	ConvertEdate, _ := time.Parse("02 January 2006", getID.Edate)
	SD := ConvertSdate.Format("2006-01-02")
	ED := ConvertEdate.Format("2006-01-02")

	// technology
	node, react, js, html := false, false, false, false
	tech := getID.Technologies
	for _ , i := range tech {
		switch i {
		case "node":
			node = true
		case "react":
			react = true
		case "js":
			js = true
		case "html5":
			html = true
		}
	}

	result := []entities.Project{getID}

	resp := map[string]interface{}{
		"Data" : entities.Data,
		"GetUpdate" : result,
		"Sdate" : SD,
		"Edate" : ED,
		"T1" : node,
		"T2" : react,
		"T3" : js,
		"T4" : html,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}


func UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := r.ParseForm()
	if err != nil{
		log.Fatal(err)
	}

	title := r.PostForm.Get("pname")
	content := r.PostForm.Get(("desc"))
	SD := r.PostForm.Get("sdate")
	ED := r.PostForm.Get("edate")
	tech := r.Form["check"]
	image := r.Context().Value("dataFile")
	img := image.(string)

	sDate, _ := time.Parse("2006-01-02", SD)
	sDateFormat := sDate.Format("02 January 2006")
	eDate, _ := time.Parse("2006-01-02", ED)
	eDateFormat := eDate.Format("02 January 2006")

	update := "UPDATE public.tb_projects SET project_name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6 WHERE id=$7"
	_, err = config.ConnDB.Exec(context.Background(), update, title, sDateFormat, eDateFormat, content, tech, img, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}


func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	deleteID := "DELETE FROM tb_projects WHERE id=$1"
	_, err := config.ConnDB.Exec(context.Background(), deleteID, id) 
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}



func Detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("./views/projectDetail.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	
	// manangkap id (id, _ (tanda _ tidak ingin menampilkan eror))
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	selectID := "SELECT id, project_name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1"
	detail := entities.Project{}
	rows := config.ConnDB.QueryRow(context.Background(), selectID, id)
	err = rows.Scan(&detail.Id, &detail.Title, &detail.Sdate, &detail.Edate, &detail.Content, &detail.Technologies, &detail.Image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// duration
	ConvertSdate, _ := time.Parse("02 January 2006", detail.Sdate)
	ConvertEdate, _ := time.Parse("02 January 2006", detail.Edate)
	duration := ConvertEdate.Sub(ConvertSdate)
	var distance string
	if duration.Hours()/24 < 7 {
		distance = strconv.FormatFloat(duration.Hours()/24, 'f', 0, 64) + " Days"
	} else if duration.Hours()/24/7 < 4 {
		distance = strconv.FormatFloat(duration.Hours()/24/7, 'f', 0, 64) + " Weeks"
	} else if duration.Hours()/24/30 < 12 {
		distance = strconv.FormatFloat(duration.Hours()/24/30, 'f', 0, 64) + " Months"
	} else {
		distance = strconv.FormatFloat(duration.Hours()/24/30/12, 'f', 0, 64) + " Years"
	}

	// technology
	node, react, js, html := false, false, false, false
	tech := detail.Technologies
	for _ , i := range tech {
		switch i {
		case "node":
			node = true
		case "react":
			react = true
		case "js":
			js = true
		case "html5":
			html = true
		}
	}


	resp := map[string]interface{}{
		"Data" : entities.Data,
		"Projects" : detail,
		"Duration" : distance,
		"T1" : node,
		"T2" : react,
		"T3" : js,
		"T4" : html,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}