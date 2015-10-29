package goldie

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"html/template"
	"log"
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

func bindArgument(t reflect.Type, values url.Values) reflect.Value {
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
	case reflect.Struct:
		return bindStruct(t, values)
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

func newHandler(action Action) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		actionType := reflect.TypeOf(action)

		args := []reflect.Value{}

		for i := 0; i < actionType.NumIn(); i++ {
			args = append(args, bindArgument(actionType.In(i), parseParameters(req)))
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
