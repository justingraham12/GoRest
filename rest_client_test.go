package GoRest
import (
    "testing"
    "fmt"
    "strings"
    "net/http/httptest"
    "net/http"
    "reflect"
)

const uri = "http://www.google.com/"

var client RestClient

func TestMakeClient(t *testing.T) {
    client = MakeClient(uri)
    if strings.Trim(uri, "/") != client.url {
        t.Error(equalsMsg(uri, client.url))
    }
}

func TestGetURL(t *testing.T) {
    if client.url != client.GetURL() {
        t.Error(equalsMsg(client.url, client.GetURL()))
    }
}

func TestGetAccept(t *testing.T) {
    if expected := client.GetAccept(); expected != ApplicationJSON {
        t.Error(equalsMsg(expected, ApplicationJSON))
    }
}

func TestPath(t *testing.T) {
    newClient := client.Path("1", "/2", "/3/")
    if expected := fmt.Sprintf("%s1/2/3", uri); expected != newClient.url {
        t.Error(equalsMsg(expected, newClient.url))
    }
}

func TestHeader(t *testing.T) {
    newClient := client.Header("key", "value")
    expected := map[string]string{"key":"value"}
    actual := newClient.GetHeaders()
    if !reflect.DeepEqual(expected, actual) {
        t.Error(equalsMsg(expected, actual))
    }
}

func TestGetJSON(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Add("Content-Type", ApplicationJSON.String())
        fmt.Fprintln(w, `{"name":"test"}`)
    }))
    defer server.Close()

    response := new(TestResponse)
    err := MakeClient(server.URL).Get(response)
    if err != nil { t.Error(err); return }

    if expected := response.Name; expected != "test" {
        t.Error(equalsMsg(expected, "test"))
    }
}

func TestGetXML(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Add("Content-Type", ApplicationXML.String())
        fmt.Fprintln(w,
        `<?xml version="1.0" encoding="UTF-8" ?>
        <Response>
            <name>test</name>
        </Response>`)
    }))
    defer server.Close()

    response := new(TestResponse)
    err := MakeClient(server.URL).
        Accept(ApplicationXML).
        Get(response)
    if err != nil { t.Error(err); return }

    if expected := response.Name; expected != "test" {
        t.Error(equalsMsg(expected, "test"))
    }
}

func TestInvalidContentType(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Add("Content-Type", ApplicationJSON.String())
        fmt.Fprintln(w, `{"name":"test"}`)
    }))
    defer server.Close()

    response := new(TestResponse)
    err := MakeClient(server.URL).
        Accept(ApplicationXML). // Set accept XML
        Get(response)
    if err == nil { t.Error("Expected Error due to Response Content-Type not matching Request Accept type") }
}

// ===================================================================
//                             Test Utils
// ===================================================================

// Builds the error message for an equals query
//
// Output: Expected [exp_val] to equal [act_val]
func equalsMsg(expected, actual interface{}) string {
    return fmt.Sprintf("Expected [%s] to equal [%s]", expected, actual)
}

// A fake response object that will be used to unmarshal test data
type TestResponse struct {
    Name string `json:"name" xml:"name"`
}

func (tr TestResponse) String() string {
    return fmt.Sprintf("TestResponse{ Name=%s }", tr.Name)
}