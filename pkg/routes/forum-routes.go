package routes

import (
	"net/http"

	"github.com/appak21/forum/pkg/controllers"
)

func RegisterForumRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./ui/css"))))
	mux.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./ui/img"))))
	mux.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./ui/js"))))

	mux.Handle("/", controllers.AppHandler(controllers.Home))
	mux.Handle("/profile", controllers.AppHandler(controllers.Profile))
	mux.Handle("/accounts/login", controllers.AppHandler(controllers.Signin))
	mux.Handle("/accounts/logout", controllers.AppHandler(controllers.Signout))
	mux.Handle("/accounts/register", controllers.AppHandler(controllers.Signup))

	mux.Handle("/create/post", controllers.AppHandler(controllers.CreatePost))
	mux.Handle("/create/comment", controllers.AppHandler(controllers.CreateComment))
	mux.Handle("/vote/post", controllers.AppHandler(controllers.VotePost))
	mux.Handle("/vote/comment", controllers.AppHandler(controllers.VoteComment))

	mux.Handle("/posts/filter", controllers.AppHandler(controllers.GetPosts))
	mux.Handle("/posts", controllers.AppHandler(controllers.GetPostByID))
	mux.Handle("/posts/myposts", controllers.AppHandler(controllers.GetPostsCreated))
	mux.Handle("/posts/mylikes", controllers.AppHandler(controllers.GetPostsLiked))

	return mux
}
