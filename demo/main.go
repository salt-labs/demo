package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/go-playground/validator.v9"

	// Go SQL driver
	_ "github.com/go-sql-driver/mysql"

	// The project is vulnerable because it uses text/template instead of html/template
	"text/template"

	// This project also uses intentionally vulnerable packages to show the difference
	// with static analysis vs variant analysis
	"github.com/go-gitea/gitea/modules/markup"
	"github.com/gophish/gophish/config"
	"golang.org/x/crypto/md4"

	// Custom module urls
	todos "demo/todos"
)

type User struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
}

var validConfig = []byte(`{
	"admin_server": {
		"listen_url": "127.0.0.1:3333",
		"use_tls": true,
		"cert_path": "gophish_admin.crt",
		"key_path": "gophish_admin.key"
	},
	"phish_server": {
		"listen_url": "0.0.0.0:8080",
		"use_tls": false,
		"cert_path": "example.crt",
		"key_path": "example.key"
	},
	"db_name": "sqlite3",
	"db_path": "gophish.db",
	"migrations_prefix": "db/db_",
	"contact_address": ""
}`)

// duplicateIf shows two conditionals in an if-else chain which are identical
func duplicateIf(msg string) {

	fmt.Println("")
	fmt.Println("Vulnerability: Duplicate If conditionals")
	fmt.Println("CWE-561 - Dead Code")

	if msg == "start" {

		fmt.Printf("\tStart conditional matched\n")

	} else if msg == "start" {

		fmt.Printf("\tStop conditional matched\n")

	} else {

		panic("Message not understood.")
	}

}

// duplicateSwitch shows two switch statements that are identical
func duplicateSwitch(msg string) {

	fmt.Println("")
	fmt.Println("Vulnerability: Duplicate Switch")
	fmt.Println("CWE-561 - Dead Code")

	switch {

	case msg == "start":

		fmt.Printf("\tStart conditional matched\n")

	case msg == "start":

		fmt.Printf("\tStop conditional matched\n")

	default:

		panic("Message not understood.")

	}

}

// inconsistentLoop demonstrates if the variable is incremented but checked against a lower bound,
// or decremented but checked against an upper bound, then the loop will usually either terminate
// immediately and never execute its body, or it will keep iterating indefinitely.
func inconsistentLoop(a []int, lower int, upper int) {

	fmt.Println("")
	fmt.Println("Vulnerability: Inconsistent Loop")
	fmt.Println("CWE-835 - Infinite Loop")

	// zero out everything below index `lower`
	for i := lower - 1; i >= 0; i-- {

		a[i] = 0

	}

	// zero out everything above index `upper`
	for i := upper + 1; i < len(a); i++ {

		a[i] = 0

	}

}

// commandInjection demonstrates why allowing user specified input
// to be executed is a bad idea
func commandInjection(w http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {

		http.Error(w, "404 not found.", http.StatusNotFound)
		return

	}

	cmdName := req.URL.Query()["cmd"][0]

	cmd := exec.Command(cmdName)

	cmd.Run()

}

// vulnValidation is an example function vulnerable due to lack of input validation
func vulnValidation() {

	fmt.Println("")
	fmt.Println("Vulnerability: Input Validation")

	v := validator.New()

	a := User{
		Email: "a", // Input an incorrect email address and lack of Name
	}

	err := v.Struct(a)

	for _, e := range err.(validator.ValidationErrors) {
		fmt.Printf("\t%v\n", e)
	}

}

// vulnSanitization is an example function vulnerable due to lack of input sanitization
func vulnSanitization() {

	fmt.Println("")
	fmt.Println("Vulnerability: Sanitization")

	input := "example<script>alert('Injected!');</script>@domain.com"

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	fmt.Printf("\tPattern: %v\n", re.String())
	fmt.Printf("\tEmail: %v :%v\n", input, re.MatchString(input))

}

// vulnInjection is an example function that is vulnerable to SQL injection attack
func vulnInjection(username string, password string) (bool, error) {

	fmt.Println("")
	fmt.Println("Vulnerability: Injection")

	db, err := sql.Open("mysql", "root:letmein@tcp(db:3306)/SocialMediaApp")

	if err != nil {
		return false, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE username ='" + username + "' AND password = '" + password + "';")
	//rows, err := db.Query("SELECT * FROM foo WHERE bar = ? AND baz = ?;", username, password)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	if rows.Next() == false {
		return false, nil
	}

	return true, nil

}

// vulnXSS is an example function that is vulnerable to cross site scripting
func vulnXSS() {

	fmt.Println("")
	fmt.Println("Vulnerability: Cross Site Scripting (XSS)")

	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		user := r.Form.Get("user")
		pw := r.Form.Get("password")

		log.Printf("Registering new user %s with password %s.\n", user, pw)

		data := todos.ToDoPageData{

			PageTitle: "My To Do's!",
			ToDos:     todos.GetToDos(),
		}

		tmpl.Execute(w, data)

	})

	http_port := os.Getenv("HTTP_PORT")
	if len(http_port) == 0 {
		http_port = "8080"
	}

	listen_addr := ":" + http_port

	http.ListenAndServe(listen_addr, nil)

}

// vulnLib is an example using intentional vulnerable third-party libraries
func vulnLib() {

	fmt.Println("")
	fmt.Println("Vulnerability: Libraries")

	// Example 1, md4 library
	h := md4.New()
	data := "These pretzels are making me thirsty."
	io.WriteString(h, data)
	fmt.Printf("\tMD4 is the new MD5: %x\n", h.Sum(nil))

	// Example 2, ioutil
	err := ioutil.WriteFile("config/phish-config.json", validConfig, 0644)
	conf := config.Config{}
	fmt.Printf("\tGone phishing for config file %v, or error: %v\n", conf, err)

	// Example 3, gitea
	fmt.Printf("\tIs the README.md file in one of the readme formats?: %v according to the gitea library\n", markup.IsReadmeFile("README.md"))

}

// main function
func main() {

	fmt.Println("hello world")
	fmt.Println("")

	// All error checking removed for testing.

	// Vulnerabilities

	// CWE-561
	duplicateIf("start")
	duplicateSwitch("message")

	// CWE-835
	var a []int
	inconsistentLoop(a, 2, 9)

	// CWE-78
	http.HandleFunc("/", commandInjection)

	vulnSanitization() // static analysis

	vulnInjection("username", "password") // static analysis

	vulnValidation() // static analysis

	vulnLib() // static analysis

	// CWE-312, CWE-315, CWE-359
	vulnXSS() // static analysis

	fmt.Println("")
	fmt.Println("Have a nice day!")

}
