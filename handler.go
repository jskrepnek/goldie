package goldie

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func bindStruct(t reflect.Type, values url.Values) reflect.Value {

	n := reflect.Indirect(reflect.New(t))

	decoder := schema.NewDecoder()
	err := decoder.Decode(n.Addr().Interface(), values)

	if err != nil {
		panic(err)
	}

	return n
}

func bindStructJson(t reflect.Type, values url.Values, b io.ReadCloser) reflect.Value {

	s := bindStruct(t, values)

	decoder := json.NewDecoder(b)
	err := decoder.Decode(s.Addr().Interface())

	if err != nil {
		panic(err)
	}

	return s
}

func bindArgument(t reflect.Type, req *http.Request) reflect.Value {

	values := parseParameters(req)

	// regardless of the type of request, primitives will only
	// be bound from the route variables and query string

	switch t.Kind() {
	case reflect.String:
		// will only bind directly when there's no ambiguity
		if len(values) == 1 {
			for key, _ := range values {
				return reflect.ValueOf(values.Get(key))
			}
		}
		panic("binding multiple primitives to action not implemented")
	case reflect.Int:
		// will only bind directly when there's no ambiguity
		if len(values) == 1 {
			for key, _ := range values {
				i, _ := strconv.Atoi(values.Get(key))
				return reflect.ValueOf(i)
			}
		}
		panic("binding multiple primitives to action not implemented")
	}

	contentType := req.Header.Get("Content-Type")
	contentType, _, _ = mime.ParseMediaType(contentType)

	switch t.Kind() {
	case reflect.Struct:

		// if it's a GET OR if it's a form encoded request, then everything
		// to bind is captured in values

		if req.Method == "GET" ||
			contentType == "application/x-www-form-urlencoded" {
			return bindStruct(t, values)
		} else {

			// if it's xml or json, then we need to bind from the
			// query string, route variables, and the body

			switch {
			case contentType == "application/json":
				return bindStructJson(t, values, req.Body)
			default:
				panic("not implemented")
			}
		}

	default:
		panic("in type not implemented")
	}

}

func parseParameters(r *http.Request) url.Values {
	r.ParseForm()

	values := url.Values{}

	for key, value := range mux.Vars(r) {
		values.Add(strings.ToLower(key), value)
	}
	for key, value := range r.Form {
		values.Add(strings.ToLower(key), value[0])
	}

	return values
}

func newModuleHandler(module *Module, inner http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		log.Printf("Module Handler: Executing %d Before hooks", len(module.Before))
		for _, before := range module.Before {
			before(req)
		}

		log.Println("Module Handler: Executing action handler")
		inner(rw, req)

		log.Printf("Module Handler: Executing %d After hooks", len(module.After))
		for _, after := range module.After {
			after(&Response{})
		}
	}
}

func newHandler(action Action) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		actionType := reflect.TypeOf(action)
		args := []reflect.Value{}

		for i := 0; i < actionType.NumIn(); i++ {
			args = append(args, bindArgument(actionType.In(i), req))
		}

		// invoke the action
		ret := reflect.ValueOf(action).Call(args)

		// handle the return value of the action

		log.Printf("First return value is of Kind %s", ret[0].Kind())

		switch ret[0].Kind() {
		case reflect.String:
			rw.Header().Add("Content-Type", "text/plain")
			fmt.Fprint(rw, ret[0].String())
		case reflect.Slice, reflect.Struct:
			log.Printf("Type of first return value: %s", ret[0].Type().Name())
			switch ret[0].Type().Name() {
			case "View":
				view := ret[0].Interface().(View)
				rw.Header().Add("Content-Type", "text/html")
				template, _ := template.New(view.Name).Parse(Templates[view.Name])
				template.Execute(rw, view.Model)
			default:
				rw.Header().Add("Content-Type", "application/json")
				encoder := json.NewEncoder(rw)
				encoder.Encode(ret[0].Interface())
			}
		}
	}
}
