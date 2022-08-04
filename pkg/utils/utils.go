package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	dir     = "ui/html"
	errFile = filepath.Join(dir, "error-page.html")
)

type appErr struct {
	Code    int
	Message string
}

func Render(w http.ResponseWriter, tmpl string, data interface{}) {
	file := filepath.Join(dir, tmpl)
	baseFile := filepath.Join(dir, "base.html")
	if !fileExists(file) || file != errFile && !fileExists(baseFile) {
		errRespond(w, appErr{Code: 404, Message: "Not Found"})
		return
	}
	if tmpls, err := template.ParseFiles(file); err == nil {
		if tmpls, err = tmpls.ParseGlob(baseFile); err == nil {
			tmpls.Execute(w, data)
			return
		}
	}
	errRespond(w, appErr{Code: 500, Message: "Internal Server Error"})
}

func HashPassword(pwd string) (string, error) {
	var pwdBytes = []byte(pwd)
	hashedPwd, err := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.MinCost)
	return string(hashedPwd), err
}

func DoPasswordsMatch(hashedPwd, currPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(currPwd))
	return err == nil
}

func Vars(r *http.Request) map[string]string {
	if rv := r.Context().Value(0); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}

func errRespond(w http.ResponseWriter, errData appErr) {
	fmt.Println("ErrRespond caused")
	tmpl, err := template.ParseFiles(errFile)
	if err != nil {
		http.Error(w, errData.Message, errData.Code)
		return
	}
	tmpl.Execute(w, &errData)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//When formats and returns time like "5 seconds ago", "6 months ago",
//past time format must be based on RFC3339 or "2006-01-02 15:04:05".
//If an error occurs, it returns past time without formatting
func When(past string) string {
	now := time.Now().Format(time.RFC3339)
	pastYear, err1 := strconv.Atoi(past[0:4])
	currYear, err2 := strconv.Atoi(now[0:4])
	if err1 != nil || err2 != nil {
		return past
	}
	if currYear-pastYear > 0 {
		return strconv.Itoa(currYear-pastYear) + " years ago"
	}

	pastMonth, err1 := strconv.Atoi(past[5:7])
	currMonth, err2 := strconv.Atoi(now[5:7])
	if err1 != nil || err2 != nil {
		return past
	}
	if currMonth-pastMonth > 0 {
		return strconv.Itoa(currMonth-pastMonth) + " months ago"
	}

	pastDay, err1 := strconv.Atoi(past[8:10])
	currDay, err2 := strconv.Atoi(now[8:10])
	if err1 != nil || err2 != nil {
		return past
	}
	if currDay-pastDay > 0 {
		return strconv.Itoa(currDay-pastDay) + " days ago"
	}

	pastHour, err1 := strconv.Atoi(past[11:13])
	currHour, err2 := strconv.Atoi(now[11:13])
	if err1 != nil || err2 != nil {
		return past
	}
	if currHour-pastHour > 0 {
		return strconv.Itoa(currHour-pastHour) + " hours ago"
	}

	pastMin, err1 := strconv.Atoi(past[14:16])
	currMin, err2 := strconv.Atoi(now[14:16])
	if err1 != nil || err2 != nil {
		return past
	}
	if currMin-pastMin > 0 {
		return strconv.Itoa(currMin-pastMin) + " minutes ago"
	}

	pastSec, err1 := strconv.Atoi(past[17:19])
	currSec, err2 := strconv.Atoi(now[17:19])
	if err1 != nil || err2 != nil {
		return past
	}
	if currSec-pastSec > 0 {
		return strconv.Itoa(currSec-pastSec) + " seconds ago"
	}

	return "just now"
}
