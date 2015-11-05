package goldie

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func buildRouter() *mux.Router {
	r := mux.NewRouter()

	for _, module := range modules {
		log.Println("buildRouter: Registering GET module routes")
		for route, action := range module.Get {
			path := fmt.Sprintf("%s%s", module.Path, route.(string))
			log.Printf("buildRouter: Registering GET %s", path)
			r.HandleFunc(path, newModuleHandler(&module, newHandler(action))).Methods("GET")
		}
		log.Println("buildRouter: Registering POST module routes")
		for route, action := range module.Post {
			path := fmt.Sprintf("%s%s", module.Path, route.(string))
			log.Printf("buildRouter: Registering POST %s", path)
			r.HandleFunc(path, newModuleHandler(&module, newHandler(action))).Methods("POST")
		}
		log.Println("buildRouter: Registering PUT module routes")
		for route, action := range module.Put {
			path := fmt.Sprintf("%s%s", module.Path, route.(string))
			log.Printf("buildRouter: Registering PUT %s", path)
			r.HandleFunc(path, newModuleHandler(&module, newHandler(action))).Methods("PUT")
		}
	}

	for route, action := range Get {
		r.HandleFunc(route.(string), newHandler(action)).Methods("GET")
	}

	for route, action := range Post {
		r.HandleFunc(route.(string), newHandler(action)).Methods("POST")
	}

	for route, action := range Put {
		r.HandleFunc(route.(string), newHandler(action)).Methods("PUT")
	}

	return r
}

func Serve() {
	http.ListenAndServe(":8080", buildRouter())
}
