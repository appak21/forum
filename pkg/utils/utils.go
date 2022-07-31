package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

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
