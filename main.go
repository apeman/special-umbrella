package main
import (
	"html/template"
	"net/http"
  	"log"
	
    "github.com/julienschmidt/httprouter"
)

const PORT = ":5678"
var tmpl = template.Must(template.ParseGlob("templates/*.html"))


func main() {
	HandleRoutes()
}


func HandleRoutes() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/dash", Dash)
	router.GET("/allf", AllFiles)
	router.GET("/del/:photoid", DeletePhoto)

	router.GET("/new", NewPhotoGet)
	router.POST("/upload", NewPhotoPost)
	router.GET("/edit/:postid", EditPost)
	router.POST("/edit/:postid", EditPost)
	router.GET("/view/:postid", ViewPost)
	router.POST("/delete/:postid", DeletePost)

//--------Login--------//
	router.GET("/register", Register)
	router.POST("/register", Register)
	router.GET("/login", loginHandler)
	router.POST("/login", loginHandler)
	router.POST("/logout", logoutHandler)

//--------FileServer--------//
	router.NotFound = http.FileServer(http.Dir(""))
	router.ServeFiles("/img/*filepath", http.Dir("uploads"))
	router.ServeFiles("/static/*filepath", http.Dir("static"))

//--------Server--------//
	log.Println("Starting erver on ", PORT)
	err := http.ListenAndServe(PORT, router)
//err := http.ListenAndServe(GetPort(), router)
 	if err != nil {
		log.Fatal("error starting http server : ", router)
 	}

}
