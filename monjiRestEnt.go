package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

//------------------------------------------------------------------------------------------------
// Data model objects and persistence mocks:
//--

type Article struct {
	ID          string `json:"id"`      // autoincremental id
	SrcType     string `json:"srcType"` // ArpltnInforInqireSvc
	Dt          string `json:"dt"`
	SidoName    string `json:"sidoName"`
	StationName string `json:"stationName"`
	KhaiGrade   int16  `json:"khaiGrade"`
	Pm10Grade   int16  `json:"pm10Grade"`
	Pm25Grade   int16  `json:"pm25Grade"`
}

// Article fixture data
var articles = []*Article{
	{ID: "1", SrcType: "ArpltnInforInqireSvc", Dt: "2019070116426", SidoName: "경기", StationName: "가평", KhaiGrade: 1, Pm10Grade: 2, Pm25Grade: 1},
	{ID: "2", SrcType: "ArpltnInforInqireSvc", Dt: "2019070116428", SidoName: "경기", StationName: "성남", KhaiGrade: 2, Pm10Grade: 1, Pm25Grade: 1},
	{ID: "3", SrcType: "ArpltnInforInqireSvc", Dt: "2019070116428", SidoName: "kyunggi", StationName: "area", KhaiGrade: 2, Pm10Grade: 1, Pm25Grade: 1},
}

func ListArticles(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewArticleListResponse(articles)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// ArticleCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var article *Article
		var err error
		var sidoName = "경기"

		if stationName := chi.URLParam(r, "stationName"); stationName != "" {

			article, err = dbGetArticle(sidoName, stationName)

		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		log.Println(r.Context())

		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateArticle persists the posted Article and returns it
// back to the client as an acknowledgement.
// func CreateArticle(w http.ResponseWriter, r *http.Request) {
// 	data := &ArticleRequest{}
// 	if err := render.Bind(r, data); err != nil {
// 		render.Render(w, r, ErrInvalidRequest(err))
// 		return
// 	}

// 	article := data.Article
// 	dbNewArticle(article)

// 	render.Status(r, http.StatusCreated)
// 	render.Render(w, r, NewArticleResponse(article))
// }

func GetArticle(w http.ResponseWriter, r *http.Request) {

	// debug contexts
	//log.Println(r.Context())
	//log.Println(r.Context().Value("StationName"))

	article := r.Context().Value("article").(*Article)

	log.Println("stop here!")
	log.Println(article)

	if err := render.Render(w, r, NewArticleResponse(article)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// type ArticleRequest struct {
// 	*Article
// 	ProtectedID string `json:"id"` // override 'id' json to have more control
// }

// func (a *ArticleRequest) Bind(r *http.Request) error {
// 	// a.Article is nil if no Article fields are sent in the request. Return an
// 	// error to avoid a nil pointer dereference.
// 	if a.Article == nil {
// 		return errors.New("missing required Article fields.")
// 	}

// 	// a.User is nil if no Userpayload fields are sent in the request. In this app
// 	// this won't cause a panic, but checks in this Bind method may be required if
// 	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

// 	// just a post-process after a decode..
// 	a.ProtectedID = "" // unset the protected ID
// 	//a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
// 	return nil
// }

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// ArticleResponse is the response payload for the Article data model.
// See NOTE above in ArticleRequest as well.
//
// In the ArticleResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type ArticleResponse struct {
	*Article
	Elapsed int64 `json:"elapsed"`
}

func NewArticleResponse(article *Article) *ArticleResponse {
	resp := &ArticleResponse{Article: article}

	return resp
}

func (rd *ArticleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func NewArticleListResponse(articles []*Article) []render.Renderer {
	list := []render.Renderer{}
	for _, article := range articles {
		list = append(list, NewArticleResponse(article))
	}
	return list
}

// func dbNewArticle(article *Article) (string, error) {
// 	article.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
// 	articles = append(articles, article)
// 	return article.ID, nil
// }

// func dbGetArticle(StationName string) (*Article, error) {
// 	for _, a := range articles {
// 		if a.StationName == StationName {
// 			return a, nil
// 		}
// 	}
// 	return nil, errors.New("article not found.")
// }

func dbGetArticle(SidoName string, StationName string) (*Article, error) {
	for _, a := range articles {
		if a.SidoName == SidoName && a.StationName == StationName {
			return a, nil
		}
	}
	return nil, errors.New("article not found.")
}
