package GoRest

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	u "net/url"
	"strings"
)

type RestClient struct {
	client      *http.Client
	url         string
	accept      MediaType
	contentType MediaType
	headers     map[string]string
	query       map[string]string
	cookies     []*http.Cookie
}

// The constructor for an immutable RestClient
//
// The baseUrl argument should be a fully qualified url and will be validated on one of the HTTP methods:
// https://github.com/ | https://groups.google.com/forum/#!forum/golang-nuts
//
// By providing the baseUrl to the RestClient it can be stored in a partial state to be built upon
// in a route or handler function
func MakeClient(baseUrl string) RestClient {
	return newClient(
		&http.Client{},
		strings.Trim(baseUrl, "/"),
		ApplicationJSON,
		ApplicationJSON,
		make(map[string]string),
		make(map[string]string),
		nil)
}

// Private constructor used to provide all RestClient parameters
func newClient(client *http.Client, url string, accept MediaType, contentType MediaType, headers map[string]string,
	query map[string]string, cookies []*http.Cookie) RestClient {
	return RestClient{
		client:      client,
		url:         url,
		accept:      accept,
		contentType: contentType,
		headers:     headers,
		query:       query,
		cookies:     cookies}
}

// ===================================================================
//                        RestClient Getters
// ===================================================================

func (rc RestClient) GetURL() string {
	uri, _ := u.Parse(rc.url)
	for k, v := range rc.query {
		uri.Query().Add(k, v)
	}
	return uri.String()
}

func (rc RestClient) GetAccept() MediaType {
	return rc.accept
}

func (rc RestClient) GetContentType() MediaType {
	return rc.contentType
}

func (rc RestClient) GetHeaders() map[string]string {
	return rc.headers
}

// ===================================================================
//                     Immutable Builder Methods
// ===================================================================

func (rc RestClient) Accept(accept MediaType) RestClient {
	return newClient(rc.client, rc.url, accept, rc.contentType, rc.headers, rc.query, rc.cookies)
}

func (rc RestClient) ContentType(contentType MediaType) RestClient {
	return newClient(rc.client, rc.url, rc.accept, contentType, rc.headers, rc.query, rc.cookies)
}

func (rc RestClient) Path(path ...string) RestClient {
	newClient := newClient(rc.client, rc.url, rc.accept, rc.contentType, rc.headers, rc.query, rc.cookies)
	for _, p := range path {
		newClient.url = fmt.Sprintf("%s/%s", newClient.url, strings.Trim(p, "/"))
	}
	return newClient
}

func (rc RestClient) Query(key, value string) RestClient {
	newQuery := make(map[string]string, len(rc.query)+1)
	for k, v := range rc.query {
		newQuery[k] = v
	}
	newQuery[key] = value
	return newClient(rc.client, rc.url, rc.accept, rc.contentType, rc.headers, newQuery, rc.cookies)
}

func (rc RestClient) Header(key, value string) RestClient {
	newHeaders := make(map[string]string)
	for k, v := range rc.headers {
		newHeaders[k] = v
	}
	newHeaders[key] = value
	return newClient(rc.client, rc.url, rc.accept, rc.contentType, newHeaders, rc.query, rc.cookies)
}

func (rc RestClient) Cookie(cookie *http.Cookie) RestClient {
	return newClient(rc.client, rc.url, rc.accept, rc.contentType, rc.headers, rc.query, append(rc.cookies, cookie))
}

func (rc RestClient) Get(entity ...interface{}) (*http.Response, error) {
	return rc.request("GET", nil, entity...)
}

func (rc RestClient) Put(reqBody []byte, entity ...interface{}) (*http.Response, error) {
	return rc.request("PUT", reqBody, entity...)
}

func (rc RestClient) Post(reqBody []byte, entity ...interface{}) (*http.Response, error) {
	return rc.request("POST", reqBody, entity...)
}

func (rc RestClient) Delete(entity ...interface{}) (*http.Response, error) {
	return rc.request("DELETE", nil, entity...)
}

// Handles building and executing the http request
func (rc RestClient) request(httpReq string, reqBody []byte, resEntity ...interface{}) (*http.Response, error) {
	// Validate the URL
	uri, _ := u.Parse(rc.url)

	// Add query params
	if len(rc.query) > 0 {
		params := u.Values{}
		for k, v := range rc.query {
			params.Add(k, v)
		}
		uri.RawQuery = params.Encode()
	}

	var (
		req *http.Request
		err error
	)

	// Build the Request
	if reqBody != nil {
		req, err = http.NewRequest(httpReq, uri.String(), bytes.NewBuffer(reqBody))
	} else {
		req, err = http.NewRequest(httpReq, uri.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Add("Accept", rc.accept.String())
	req.Header.Add("Content-Type", rc.contentType.String())

	for k, v := range rc.headers {
		req.Header.Add(k, v)
	}

	// Make Request
	res, err := rc.client.Do(req)
	if err != nil {
		return res, err
	}

	return rc.response(res, resEntity...)
}

// Handles validating and reading the http response
func (rc RestClient) response(res *http.Response, resEntity ...interface{}) (*http.Response, error) {
	// Validate the response content type matches the accept type.
	// This is required to allow unmarshalling to the resEntity
	if contentType := res.Header.Get("Content-Type"); len(resEntity) != 0 &&
		!strings.Contains(strings.ToLower(contentType), strings.ToLower(rc.accept.String())) {
		return res, errors.New(fmt.Sprintf("Expected Response Content-Type [%s] to match/contain Request Accept [%s]",
			contentType, rc.accept.String()))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}

	return res, rc.unmarshal(body, resEntity...)
}

// Tries converting the body into the provided entities. This function will only return an error
// if it is unable to unmarshal to all provided entities.
func (rc RestClient) unmarshal(body []byte, resEntity ...interface{}) error {
	if len(resEntity) > 0 {
		e := []error{}

		for _, entity := range resEntity {
			if err := rc.accept.Unmarshal(body, entity); err != nil {
				e = append(e, err)
			}
		}

		if len(e) == len(resEntity) {
			return fmt.Errorf("Error during unmarshal: %v", e)
		}
	}

	return nil
}
