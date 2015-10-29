package goldie

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {

	Get["/no/params"] = func() string {
		return "hello"
	}

	Get["/one/variable/{id}"] = func(m struct{ Id int }) string {
		return strconv.Itoa(m.Id)
	}

	Get["/mixed/{name}/types/{age}"] = func(m struct {
		Name  string
		Age   int
		Num   int
		Extra int
	}) string {
		return fmt.Sprintf("Hello %s you are %d years old.  Here are %d things.", m.Name, m.Age, m.Num+m.Extra)
	}

	Get["/boolean"] = func(m struct{ Male bool }) string {
		return fmt.Sprintf("Male = %v", m.Male)
	}

	type Widget struct {
		Foo int
		Bar string
	}

	Get["/widgets"] = func() []Widget {
		var widgets = []Widget{
			{
				Foo: 10,
				Bar: "flip",
			},
			{
				Foo: 40,
				Bar: "flop",
			},
		}
		return widgets
	}

	Get["/widgets/{id}"] = func(id string) Widget {
		return Widget{
			Foo: 15,
			Bar: id,
		}
	}

	Get["/widgets/type/{id}"] = func(id int) []Widget {
		return []Widget{
			{
				Foo: id,
				Bar: "type",
			},
		}
	}

	Templates["home"] =
		`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Home</title>
		</head>
		<body>
			Welcome Home {{.Name}}!
		</body>
	</html>							
	`

	Get["/home"] = func() View {
		return View{"home", struct{ Name string }{"Joel"}}
	}

	var tests = []struct {
		name   string
		route  string
		result string
	}{
		{
			name:   "No variables",
			route:  "/no/params",
			result: "hello",
		},
		{
			name:   "int variable",
			route:  "/one/variable/34343",
			result: "34343",
		},
		{
			name:   "int and string variables",
			route:  "/mixed/bob/types/20",
			result: "Hello bob you are 20 years old.  Here are 0 things.",
		},
		{
			name:   "query parameters",
			route:  "/mixed/bob/types/20?num=300",
			result: "Hello bob you are 20 years old.  Here are 300 things.",
		},
		{
			name:   "multiple query parameters",
			route:  "/mixed/box/types/40?num=300&extra=100",
			result: "Hello box you are 40 years old.  Here are 400 things.",
		},
		{
			name:   "boolean parameter",
			route:  "/boolean?male=true",
			result: "Male = true",
		},
		{
			name:   "boolean parameter",
			route:  "/boolean?male=false",
			result: "Male = false",
		},
		{
			name:   "return a json array response",
			route:  "/widgets",
			result: "[{\"Foo\":10,\"Bar\":\"flip\"},{\"Foo\":40,\"Bar\":\"flop\"}]\n",
		},
		{
			name:   "return a json object binding directly to a single string param",
			route:  "/widgets/3434",
			result: "{\"Foo\":15,\"Bar\":\"3434\"}\n",
		},
		{
			name:   "template",
			route:  "/home",
			result: "\n\t<!DOCTYPE html>\n\t<html>\n\t\t<head>\n\t\t\t<meta charset=\"UTF-8\">\n\t\t\t<title>Home</title>\n\t\t</head>\n\t\t<body>\n\t\t\tWelcome Home Joel!\n\t\t</body>\n\t</html>\t\t\t\t\t\t\t\n\t",
		},
		{
			name:   "binding to a single int parameter",
			route:  "/widgets/type/323424",
			result: "[{\"Foo\":323424,\"Bar\":\"type\"}]\n",
		},
	}

	server := httptest.NewServer(buildRouter())
	defer server.Close()

	for _, test := range tests {
		func() {
			response, _ := http.Get(server.URL + test.route)
			defer response.Body.Close()

			content, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, test.result, string(content))
		}()
	}
}
