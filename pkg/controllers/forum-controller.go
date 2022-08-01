package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/appak21/forum/pkg/models"
	"github.com/appak21/forum/pkg/utils"
)

const TimeFormat string = time.RFC1123

type appData struct {
	User           models.User
	Posts          []models.Post
	Post           models.Post //replace
	Tags           []string
	Tag            string
	TotalPosts     int
	Path           string
	WarningMessage string
	SessionOpen    bool
}

type appError struct {
	Code    int
	Message string
}

type AppHandler func(http.ResponseWriter, *http.Request) *appError

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RECOVER ERR:", err)
			w.WriteHeader(500)
			utils.Render(w, "error-page.html", appError{500, http.StatusText(500)})
		}
	}()
	if appErr := fn(w, r); appErr != nil {
		w.WriteHeader(appErr.Code)
		appErr.Message = http.StatusText(appErr.Code)
		utils.Render(w, "error-page.html", appErr)
	}
}

func Home(w http.ResponseWriter, r *http.Request) *appError {
	if r.URL.Path != "/" {
		return &appError{Code: http.StatusNotFound}
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	posts, err := models.GetAllPosts()
	if err != nil {
		return &appError{Code: http.StatusInternalServerError}
	}
	tags, err := models.GetAllTags()
	if err != nil {
		return &appError{Code: http.StatusInternalServerError}
	}
	username, isSessionOpen := ValidSession(r)
	data := &appData{
		SessionOpen: isSessionOpen,
		User:        models.User{Username: username},
		Posts:       *posts,
		Tags:        tags,
		TotalPosts:  len(*posts),
	}
	utils.Render(w, "home-page.html", data)
	return nil
}

func Profile(w http.ResponseWriter, r *http.Request) *appError {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	data := &appData{}
	var username string
	login := r.URL.Query().Get("username")
	username, data.SessionOpen = ValidSession(r)
	if login != "" && login != username {
		username = login
	}
	user, err := models.GetUser(username)
	if err != nil {
		fmt.Println("User not found", err)
		return &appError{Code: http.StatusNotFound}
	}
	posts, err := models.GetPostsCreatedByUser(user.ID)
	if err != nil {
		return &appError{Code: 500}
	}
	for _, post := range *posts {
		user.Reputation += post.Votes.Likes - post.Votes.Dislikes
	}
	data.User = *user
	data.TotalPosts = len(*posts)
	utils.Render(w, "profile-page.html", data)
	return nil
}

func Signin(w http.ResponseWriter, r *http.Request) *appError {
	nextURL := r.FormValue("next")
	if nextURL == "" {
		nextURL = "/"
	}
	_, isSessionOpen := ValidSession(r)
	if isSessionOpen {
		http.Redirect(w, r, nextURL, http.StatusSeeOther)
		return nil
	}
	data := &appData{Path: nextURL}
	switch r.Method {
	case http.MethodGet:
		utils.Render(w, "signin-page.html", data)
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			return &appError{Code: http.StatusBadRequest}
		}
		user, err := models.GetUser(username)
		if err != nil || !utils.DoPasswordsMatch(user.Password, password) {
			data.WarningMessage = "Incorrect username or password."
			utils.Render(w, "signin-page.html", data)
			return nil
		}
		NewSessionToken(w, username)
		http.Redirect(w, r, nextURL, http.StatusSeeOther)
	default:
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	return nil
}

func Signout(w http.ResponseWriter, r *http.Request) *appError {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	_, isSessionOpen := ValidSession(r)
	if !isSessionOpen {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}
	c, _ := r.Cookie("session_token")
	sessions.Delete(c.Value)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func Signup(w http.ResponseWriter, r *http.Request) *appError {
	_, isSessionOpen := ValidSession(r)
	if isSessionOpen {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}
	data := &appData{}
	switch r.Method {
	case http.MethodGet:
		utils.Render(w, "signup-page.html", data)
	case http.MethodPost:
		time := time.Now().Format(TimeFormat)
		user := models.User{
			Username:  r.FormValue("username"),
			Email:     r.FormValue("email"),
			Password:  r.FormValue("password"),
			CreatedAt: time,
		}
		confirmPwd := r.FormValue("confirm")
		//----------
		if user.Username == "" || user.Email == "" || user.Password == "" || confirmPwd == "" {
			return &appError{Code: http.StatusBadRequest}
		}
		if len(user.Username) > 16 || len(user.Password) < 6 || user.Password != confirmPwd {
			return &appError{Code: 400}
		}
		loginExpr := "^[a-zA-Z0-9]*$"
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
		if rex, _ := regexp.Compile(loginExpr); !rex.MatchString(user.Username) {
			return &appError{Code: 400}
		}
		if _, err := mail.ParseAddress(user.Email); err != nil || !emailRegex.MatchString(user.Email) {
			return &appError{Code: 400}
		}
		//----------
		hashedPwd, err := utils.HashPassword(user.Password)
		if err != nil {
			return &appError{Code: 500}
		}
		user.Password = hashedPwd
		if err = models.CreateUser(user); err != nil {
			data.WarningMessage = "The username or email already exists" //be more clear
			utils.Render(w, "signup-page.html", data)
			return nil
		}
		NewSessionToken(w, user.Username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		return &appError{Code: 405}
	}
	return nil
}

func CreatePost(w http.ResponseWriter, r *http.Request) *appError {
	username, isSessionOpen := ValidSession(r)
	if !isSessionOpen {
		http.Redirect(w, r, "/accounts/login?next=/create/post", http.StatusFound)
		return nil
	}
	data := &appData{SessionOpen: isSessionOpen}
	switch r.Method {
	case http.MethodGet:
		utils.Render(w, "create-page.html", data)
	case http.MethodPost:
		time := time.Now().Format(TimeFormat)
		post := &models.Post{
			Title:     r.FormValue("title"),
			Text:      r.FormValue("text"),
			Tags:      strings.Split(r.FormValue("tags"), " "),
			CreatedAt: time,
		}
		fmt.Println("TAGS:", post.Tags)
		//------------
		if strings.TrimSpace(post.Title) == "" || strings.TrimSpace(post.Text) == "" {
			return &appError{Code: http.StatusBadRequest}
		}
		if utf8.RuneCountInString(post.Title) > 100 || utf8.RuneCountInString(post.Text) > 10000 {
			return &appError{Code: http.StatusBadRequest}
		}
		if len(post.Tags) > 50 {
			return &appError{Code: http.StatusBadRequest}
		}
		for _, tag := range post.Tags {
			if strings.Contains(tag, " ") || utf8.RuneCountInString(tag) > 30 || tag == "" {
				return &appError{Code: http.StatusBadRequest}
			}
		}
		//------------
		user, err := models.GetUser(username)
		if err != nil {
			return &appError{Code: 500}
		}
		post.UserID = user.ID
		post.Username = username
		if err = models.CreatePost(post); err != nil {
			return &appError{Code: 500}
		}
		if err = models.CreateTags(post.ID, post.Tags); err != nil {
			return &appError{Code: 500}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	return nil
}

func CreateComment(w http.ResponseWriter, r *http.Request) *appError {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}
	username, isSessionOpen := ValidSession(r)
	if !isSessionOpen {
		http.Redirect(w, r, "/accounts/login/?next=/create/comment", http.StatusFound)
		return nil
	}
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return &appError{Code: http.StatusNotFound}
	}
	user, err := models.GetUser(username)
	if err != nil {
		return &appError{Code: http.StatusInternalServerError}
	}
	time := time.Now().Format(TimeFormat)
	comment := &models.Comment{
		PostID:    int64(postID),
		UserID:    user.ID,
		Username:  username,
		Text:      r.FormValue("text"),
		CreatedAt: time,
	}
	//-------------
	if strings.TrimSpace(comment.Text) == "" || utf8.RuneCountInString(comment.Text) > 200 {
		return &appError{Code: http.StatusBadRequest}
	}
	//-------------
	if err = models.CreateComment(*comment); err != nil {
		return &appError{Code: 500}
	}
	next := fmt.Sprintf("/posts?id=%v", postID)
	http.Redirect(w, r, next, http.StatusSeeOther)
	return nil
}

func VotePost(w http.ResponseWriter, r *http.Request) *appError {
	username, isSessionOpen := ValidSession(r)
	if !isSessionOpen {
		http.Redirect(w, r, "/accounts/login/?next=/vote/post", http.StatusFound) //get vote/post?...
		return nil
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return &appError{Code: http.StatusBadRequest}
	}
	vote, err := strconv.Atoi(r.URL.Query().Get("vote"))
	if err != nil {
		return &appError{Code: http.StatusBadRequest}
	}
	if vote != 1 && vote != -1 {
		return &appError{Code: http.StatusBadRequest}
	}
	user, err := models.GetUser(username)
	if err != nil {
		return &appError{Code: http.StatusBadRequest}
	}
	if _, err = models.GetPostById(int64(postID)); err != nil {
		return &appError{Code: http.StatusBadRequest}
	}

	if err = models.VotePost(user.ID, int64(postID), vote); err != nil {
		return &appError{Code: 500}
	}
	next := fmt.Sprintf("/posts/%v", postID)
	http.Redirect(w, r, next, http.StatusSeeOther)
	return nil
}

func VoteComment(w http.ResponseWriter, r *http.Request) *appError {
	username, isSessionOpen := ValidSession(r)
	if !isSessionOpen {
		http.Redirect(w, r, "/accounts/login/?next=/vote/comment", http.StatusFound)
		return nil
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		return &appError{Code: http.StatusMethodNotAllowed}
	}

	cmtID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return &appError{Code: http.StatusBadRequest}
	}
	vote, err := strconv.Atoi(r.URL.Query().Get("vote"))
	if err != nil {
		return &appError{Code: http.StatusBadRequest}
	}
	if vote != 1 && vote != -1 {
		return &appError{Code: http.StatusBadRequest}
	}
	user, err := models.GetUser(username)
	if err != nil {
		return &appError{Code: http.StatusBadRequest}
	}
	comment, err := models.GetCommentByID(int64(cmtID))
	if err != sql.ErrNoRows {
		return &appError{Code: http.StatusBadRequest}
	} else if err != nil {
		return &appError{Code: 500}
	}
	if err = models.VoteComment(user.ID, int64(cmtID), vote); err != nil {
		return &appError{Code: 500}
	}
	next := fmt.Sprintf("/posts/%v", comment.PostID)
	http.Redirect(w, r, next, http.StatusSeeOther)
	return nil
}

func GetPosts(w http.ResponseWriter, r *http.Request) *appError { //FILTER FUNC
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	_, isSessionOpen := ValidSession(r)
	var data = &appData{SessionOpen: isSessionOpen}
	tag := r.FormValue("tag")
	tag = strings.TrimPrefix(tag, "#")
	if tag != "" {
		posts, err := models.GetPostsByTag(tag)
		if err == sql.ErrNoRows { //
			return &appError{Code: http.StatusNotFound}
		} else if err != nil {
			return &appError{Code: 500}
		}
		data.Posts = *posts
		data.TotalPosts = len(*posts)
	} else {
		posts, err := models.GetAllPosts()
		if err != nil {
			return &appError{Code: 500}
		}
		data.Posts = *posts
		data.TotalPosts = len(*posts)
	}
	data.Tag = tag
	allTags, err := models.GetAllTags()
	if err != nil {
		return &appError{Code: 500}
	}
	data.Tags = allTags
	utils.Render(w, "home-page.html", data)
	return nil
}

func GetPostByID(w http.ResponseWriter, r *http.Request) *appError {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	val := r.URL.Query().Get("id")
	if val == "" {
		return GetPosts(w, r)
	}
	postID, err := strconv.ParseInt(val, 0, 0)
	if err != nil {
		return &appError{Code: http.StatusNotFound}
	}
	post, err := models.GetPostById(postID)
	if err != nil {
		return &appError{Code: 404}
	}
	username, isSessionOpen := ValidSession(r)
	data := &appData{
		SessionOpen: isSessionOpen,
		User:        models.User{Username: username},
		// Posts:       []models.Post{*post}, //replace with this
		Post: *post,
	}
	utils.Render(w, "post-page.html", data)
	return nil
}

func GetPostsCreated(w http.ResponseWriter, r *http.Request) *appError {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	login := r.URL.Query().Get("login")
	username, isSessionOpen := ValidSession(r)
	if login == "" && !isSessionOpen {
		http.Redirect(w, r, "/accounts/login/?next=/posts/myposts", http.StatusSeeOther)
		return nil
	}
	if login == "" {
		login = username
	}
	user, err := models.GetUser(login)
	if err != nil {
		return &appError{Code: http.StatusNotFound}
	}
	posts, err := models.GetPostsCreatedByUser(user.ID)
	if err != nil {
		return &appError{Code: http.StatusInternalServerError}
	}
	tags, err := models.GetAllTags()
	if err != nil {
		return &appError{Code: 500}
	}
	utils.Render(w, "home-page.html", appData{
		SessionOpen: isSessionOpen,
		User:        *user,
		Posts:       *posts,
		Tags:        tags,
		TotalPosts:  len(*posts),
	})
	return nil
}

func GetPostsLiked(w http.ResponseWriter, r *http.Request) *appError {
	username, isSessionOpen := ValidSession(r)
	if !isSessionOpen {
		http.Redirect(w, r, "/accounts/login/?next=/posts/mylikes", http.StatusFound)
		return nil
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		return &appError{Code: http.StatusMethodNotAllowed}
	}
	user, _ := models.GetUser(username)
	posts, err := models.GetPostsVotedByUser(user.ID, 1)
	if err != nil {
		return &appError{Code: 500}
	}
	tags, err := models.GetAllTags()
	if err != nil {
		return &appError{Code: 500}
	}
	data := &appData{
		SessionOpen: isSessionOpen,
		Posts:       *posts,
		User:        *user,
		Tags:        tags,
		TotalPosts:  len(*posts),
	}
	utils.Render(w, "home-page.html", data)
	return nil
}
