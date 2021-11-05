package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"myWeb/configs"
	"net/http"
	"text/template"
)

type Article struct {
	Id                    uint16
	Title, Path, Password string
}

var posts []Article
var showPost Article

var config = loadConfig()

func main() {
	handlFunc()
}

func loadConfig() *configs.Config {
	config := configs.NewConfig()

	yamFile, err := ioutil.ReadFile("./configs/app.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamFile, config)
	if err != nil {
		panic(err)
	}

	return config
}

func reverse(s []Article) []Article {
	a := make([]Article, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html",
		"templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s@(%s%s)/golang",
		config.User, config.DbHost, config.DbPort))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM articles")
	if err != nil {
		panic(err)
	}

	posts = []Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Path, &post.Password)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	err = t.ExecuteTemplate(w, "index", reverse(posts))
	if err != nil {
		panic(err)
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html",
		"templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	err = t.ExecuteTemplate(w, "create", nil)
	if err != nil {
		panic(err)
	}
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	path := r.FormValue("path")
	password := r.FormValue("password")

	if password != config.Password {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	} else if title == "" || path == "" || password == "" {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	} else {
		db, err := sql.Open("mysql", fmt.Sprintf("%s@(%s%s)/golang",
			config.User, config.DbHost, config.DbPort))
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO articles (title, path, password) VALUES('%s', '%s', '%s')",
			title, path, password))
		if err != nil {
			panic(err)
		}

		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html",
		"templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s@(%s%s)/golang",
		config.User, config.DbHost, config.DbPort))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE id = '%s'",
		vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Path, &post.Password)
		if err != nil {
			panic(err)
		}

		showPost = post
	}

	err = t.ExecuteTemplate(w, "show", showPost)
	if err != nil {
		panic(err)
	}

}

func show_error(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/error.html",
		"templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	err = t.ExecuteTemplate(w, "error", nil)
	if err != nil {
		panic(err)
	}
}

func handlFunc() {
	config := loadConfig()

	rtr := mux.NewRouter()

	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	rtr.HandleFunc("/error", show_error).Methods("GET")

	http.Handle("/", rtr)

	http.ListenAndServe(config.Port, nil)
	//if err != nil {
	//	panic(err)
	//}
}
