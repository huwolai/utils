package route

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"time"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type IHandler interface  {

	Handler(inner http.Handler, name string) http.Handler
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func NewRouterWithHandle(routes []Route,handles []IHandler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc

		rout := router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name)
		if handles !=nil {
			for _,handl := range handles {
				handler = handl.Handler(handler,route.Name)
				rout = rout.Handler(handler)
			}
		}

	}
	return router
}

func NewRouter(routes []Route) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
	}
	return router
}
