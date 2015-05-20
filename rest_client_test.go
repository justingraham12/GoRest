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

func TestGetAccept(t *testing.T) {
	assert.Equal(t, client.GetAccept(), ApplicationJSON)
}

func TestPath(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%s1/2/3", uri), client.Path("1", "/2", "/3/").url)
}

func TestHeader(t *testing.T) {
	assert.Equal(t, map[string]string{"key": "value"}, client.Header("key", "value").GetHeaders())
}

func TestGetJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", ApplicationJSON.String())
		fmt.Fprintln(w, `{"name":"test"}`)
	}))
	defer server.Close()

	response := new(u.TestResponse)
	_, err := MakeClient(server.URL).Get(response)
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

	response := new(u.TestResponse)
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
		Get(new(u.TestResponse))
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

	response := new(u.TestResponse)
	_, err := MakeClient(server.URL).
		Accept(ApplicationJSON).
		Query("key", "value").
		Get(response)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "passed", response.Name)
}
