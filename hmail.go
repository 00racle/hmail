package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"fmt"
	"time"
	"strings"
	"html/template"
)

var db *sql.DB
var err error

func SuccessPage(res http.ResponseWriter, req *http.Request){
	if req.Method != "POST"{
		http.ServeFile(res, req, "success.html")
		return
	}
}

func Mpage(res http.ResponseWriter, req *http.Request){
	ipport := req.RemoteAddr
	ipaddr := strings.Split(ipport, ":")[0]
	if req.Method != "POST"{
		http.ServeFile(res, req, "nsite/index.html")
		fmt.Println(req.Method, ipaddr)
		_, err = db.Exec("INSERT INTO access(ip) VALUES(?)",ipaddr)
		if err != nil {
			panic(err.Error())
			}
		return
	}
	redirectTarget := "/success"
	t := time.Now().UTC()
	t = t.In(time.FixedZone("KST", 9*60*60))

	username := req.FormValue("username")
	password := req.FormValue("password")

	// 패스워드 일방향 암화화
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}

	accesstime := t.Format("2006-01-02 15:04:05")

	_, err = db.Exec("INSERT INTO users(ip, username, password) VALUES(?, ?, ?)",ipaddr, username[:2], hashedPassword)
	if err != nil {
		panic(err.Error())
	}
	http.Redirect(res, req, redirectTarget, 302)
	fmt.Println("접속시간: [",accesstime,"]","접속IP: [",ipaddr,"]","ID: [",username[:2],"]","비밀번호: [",hashedPassword,"]")
}
func ResultPage(res http.ResponseWriter, req *http.Request){
	//t, err := template.ParseFiles("result.html")
	t, err := template.ParseFiles("liq.html")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	var a, c int
	
	db.QueryRow("SELECT COUNT(DISTINCT ip) from access").Scan(&a)
	db.QueryRow("SELECT COUNT(DISTINCT ip) from users").Scan(&c)

	//items := map[string]int
	var items map[string]int
	items = make(map[string]int)
	
	items["total"] = 500
	items["A_count"] = a
	items["C_count"] = c


	if err := t.Execute(res, items); err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}

func main(){
	db, err = sql.Open("mysql", "root:8282op82@#@tcp(localhost:3306)/dchk")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.Handle("/", http.FileServer(http.Dir("nsite/")))

	http.HandleFunc("/result", ResultPage)
	http.HandleFunc("/success", SuccessPage)
	http.HandleFunc("/login", Mpage)
	http.ListenAndServe(":8632", nil)
}
