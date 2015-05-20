package GoRest

import (
	u "./utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const uri = "http://www.google.com/"

var client RestClient

func TestMakeClient(t *testing.T) {
	client = MakeClient(uri)
	assert.Equal(t, strings.Trim(uri, "/"), client.url)
}

func TestGetURL(t *testing.T) {
	assert.Equal(t, client.url, client.GetURL())
}

func TestGetURLWithQuery(t *testing.T) {
	cli := MakeClient(uri).Query("key", "value")
	assert.Equal(t, "http://www.google.com?key=value", cli.GetURL())
}

func TestGetAccept(t *testing.T) {
	assert.Equal(t, client.GetAccept(), ApplicationJSON)
}

func TestContentType(t *testing.T) {
	client.ContentType(ApplicationJSON)
	assert.Equal(t, client.GetContentType(), ApplicationJSON)
}

func TestPath(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%s1/2/3", uri), client.Path("1", "/2", "/3/").url)
}

func TestHeader(t *testing.T) {
	assert.Equal(t,
		map[string]string{"key": "value", "key2": "value2"},
		client.Header("key", "value").Header("key2", "value2").GetHeaders())
}

func TestGetJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", ApplicationJSON.String())
		fmt.Fprintln(w, `{"name":"test"}`)
	}))
	defer server.Close()

	response := new(u.TestResponse1)
	_, err := MakeClient(server.URL).Get(response)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "test", response.Name)
}

func TestPostJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", ApplicationJSON.String())
		fmt.Fprintln(w, `{"name":"test"}`)
	}))
	defer server.Close()

	response := new(u.TestResponse1)
	_, err := MakeClient(server.URL).Post([]byte{}, response)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "test", response.Name)
}

func TestGetXML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", ApplicationXML.String())
		fmt.Fprintln(w, `<Response><name>test</name></Response>`)
	}))
	defer server.Close()

	response := new(u.TestResponse1)
	_, err := MakeClient(server.URL).
		Accept(ApplicationXML).
		Get(response)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "test", response.Name)
}

func TestInvalidContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", ApplicationJSON.String()) // Set Content-Type JSON
		fmt.Fprintln(w, `{"name":"test"}`)
	}))
	defer server.Close()

	_, err := MakeClient(server.URL).
		Accept(ApplicationXML). // Set accept type XML
		Get(new(u.TestResponse1))
	assert.NotNil(t, err, "Expected Error due to Accept type not matching Content-Type")
}

func TestQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", ApplicationJSON.String())
		if query := r.URL.Query(); query.Get("key") == "value" {
			fmt.Fprintln(w, `{"name":"passed"}`)
		} else {
			fmt.Fprintln(w, `{"name":"failed"}`)
		}
	}))
	defer server.Close()

	response := new(u.TestResponse1)
	_, err := MakeClient(server.URL).Query("key", "value").Query("key2", "value2").Get(response)

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "passed", response.Name)
}

func TestMultipleEntity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", ApplicationJSON.String())
		fmt.Fprintln(w, `{"name":"test", "age":25}`)
	}))
	defer server.Close()

	response1 := new(u.TestResponse1)
	response2 := new(u.TestResponse2)
	_, err := MakeClient(server.URL).Get(response1, response2)

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "test", response1.Name)
	assert.Equal(t, 25, response2.Age)
}

func TestArrayEntity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", ApplicationJSON.String())
		fmt.Fprintln(w, `[{"name":"test", "age":25}]`)
	}))
	defer server.Close()

	response1 := new(u.TestResponseArray)
	response2 := new(u.TestResponse2)
	_, err := MakeClient(server.URL).Get(response2, response1)

	if err != nil {
		t.Error(err)
		return
	}

	assert.Len(t, *response1, 1)
	assert.Equal(t, "test", (*response1)[0].Name)
	assert.Equal(t, 0, response2.Age)
}

func TestMultipleEntityFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", ApplicationJSON.String())
		fmt.Fprintln(w, `[{"hello":"test", "world":"25"}]`)
	}))
	defer server.Close()

	response1 := new(u.TestResponse1)
	response2 := new(u.TestResponse2)
	_, err := MakeClient(server.URL).Get(response1, response2)

	assert.NotNil(t, err, "Expected Unmarshal Error")
}
