package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	// _ "net/http/pprof"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	// RESTy routes for "articles" resource
	r.Route("/locations", func(r chi.Router) {
		r.With(paginate).Get("/", ListArticles)
		//r.Get("/search", SearchArticles) // GET /articles/search

		r.Route("/{stationName}", func(r chi.Router) {
			r.Use(ArticleCtx)      // Load the *Article on the request context
			r.Get("/", GetArticle) // GET /articles/123

		})

	})

	//-------------------------------------------------------------------------------------------------------------------------------------
	//mysql connection
	//-------------------------------------------------------------------------------------------------------------------------------------

	// db, err := sql.Open("mysql", "jonas:@tcp(127.0.0.1:3306)/test")
	// //db, err := sql.Open("mysql", "admin:@tcp(10.178.230.98:3306)/monji")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// // select one recent data and print

	// result, err := db.Exec("SELECT * FROM monji limit 1")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println("result:")
	// //fmt.Println(result)
	// log.Println(result)

	// for pprof
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	http.ListenAndServe(":3333", r)

}
