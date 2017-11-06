package httptest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

//T is a struct used of writing tests, contains request information.
type T struct {
	path      string
	method    string
	f         func(http.ResponseWriter, *http.Request)
	manualReq bool
	t         *testing.T
	req       *http.Request
	params    *url.Values
	postCt    string
	done      bool
	rr        *httptest.ResponseRecorder

	ba         bool
	baNamename string
	baPassword string

	cookies []*http.Cookie
}

//New is a function create a new http request, with the url and function being tested .
func New(path string, f func(http.ResponseWriter, *http.Request), t *testing.T) *T {
	r := &T{}
	r.path = path
	r.manualReq = false
	r.method = "GET"
	r.done = false
	r.params = &url.Values{}
	r.f = f
	r.t = t
	r.postCt = "nil"

	return r
}

//SetRequest is to set your custom *http.Request to the test.
//Deafult method is GET,
//Default ContentType of POST is application/x-www-form-urlencoded.
func (t *T) SetRequest(r *http.Request) *T {
	t.manualReq = true
	t.req = r
	return t
}

//SetContentType is used less except doing post.
//Default content-type of a post request is application/x-www-form-urlencoded.
func (t *T) SetContentType(contentType string) *T {
	t.postCt = contentType
	return t
}

//SetContentTypeFormUrlencoded only works for post.
func (t *T) SetContentTypeFormUrlencoded() *T {
	t.postCt = "application/x-www-form-urlencoded"
	return t
}

//SetContentTypeMultipart only works for post.
func (t *T) SetContentTypeMultipart() *T {
	t.postCt = "multipart/form-data"
	return t
}

//Post means the method is post.
//Default ContentType of POST is application/x-www-form-urlencoded.
func (t *T) Post() *T {
	t.method = "POST"
	return t
}

//Get means the method is post.
func (t *T) Get() *T {
	t.method = "GET"
	return t
}

//Put means the method is post.
func (t *T) Put() *T {
	t.req.Method = "PUT"
	return t
}

//Delete means the method is post.
func (t *T) Delete() *T {
	t.req.Method = "DELETE"
	return t
}

//Patch means the method is post.
func (t *T) Patch() *T {
	t.req.Method = "PATCH"
	return t
}

//Head means the method is post.
func (t *T) Head() *T {
	t.req.Method = "HEAD"
	return t
}

//SetBasicAuth set the request to use basic http auth.
func (t *T) SetBasicAuth(username, password string) *T {
	t.ba = true
	t.baNamename = username
	t.baPassword = password
	return t
}

//AddCookies adds cookie(s) to the request.
func (t *T) AddCookies(c ...*http.Cookie) *T {
	t.cookies = append(t.cookies, c...)
	return t
}

//AddParams add parameters to request.
//In get request, the parameters are encoded in url.Values(url)
//In other requests, the parameters are encoded in request body.
func (t *T) AddParams(k, v string) *T {
	if t.params == nil {
		t.params = &url.Values{}
	}
	t.params.Add(k, v)
	return t
}

//ResponseRecorder returns the origin *httptest.ResponseRecorder
func (t *T) ResponseRecorder() *httptest.ResponseRecorder {
	t.checkDone()
	return t.rr
}

func (t *T) checkDone() {
	if !t.done {
		t.t.Errorf("Request of [%s] to [%s] have not been done, cannot get the recoder", t.method, t.path)
	}
}

//CheckCode checks whether the response code equals the expected.
func (t *T) CheckCode(code int) *T {
	t.checkDone()
	if t.rr.Code != code {
		t.t.Errorf("Request of [%s] to [%s] test error, want response code [%d], but got [%d]", t.method, t.path, code, t.rr.Code)
	}
	return t
}

//CheckHeader checks whether the header value equals the expected.
func (t *T) CheckHeader(name, want string) *T {
	t.checkDone()
	actual := t.rr.Header().Get(name)

	if actual != want {
		t.t.Errorf("Request of [%s] to [%s] test error, want header [%s] to equal: [%s], but got: [%s]", t.method, t.path, name, want, actual)
	}
	return t
}

//BodyContains checks whether the body contains a certain string.
func (t *T) BodyContains(want string) *T {
	t.checkDone()
	if t.rr == nil || t.rr.Body == nil {
		return t
	}
	if !strings.Contains(string(t.rr.Body.Bytes()), want) {
		t.t.Errorf("Request of [%s] to [%s] test error, want content contains [%s], but got none", t.method, t.path, want)
	}
	return t
}

//Body returns the pure string(response body).
func (t *T) Body() string {
	t.checkDone()
	return string(t.rr.Body.Bytes())
}

//Do just make a request.
//CheckCode, CheckHeader, BodyContains and Body rely on Do.
//Just do it.
func (t *T) Do() *T {
	t.done = true
	paramsEncoded := t.params.Encode()
	reader := strings.NewReader(paramsEncoded)

	if !t.manualReq {
		if t.method == "GET" {
			t.req, _ = http.NewRequest(t.method, t.path, nil)
		} else {
			t.req, _ = http.NewRequest(t.method, t.path, reader)
		}
	}

	for _, v := range t.cookies {
		t.req.AddCookie(v)
	}
	if t.ba {
		t.req.SetBasicAuth(t.baNamename, t.baPassword)
	}
	if t.req.Method == "GET" {
		q := t.req.URL.Query()
		for k, v := range *t.params {
			for _, v2 := range v {
				q.Add(k, v2)
			}
		}
		t.req.URL.RawQuery = q.Encode()
	}

	if t.req.Method == "POST" && t.postCt != "nil" {
		t.req.Header.Set("Content-Type", t.postCt)
	}

	if t.req.Method == "POST" && t.postCt == "nil" {
		t.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	t.rr = httptest.NewRecorder()
	handler := http.HandlerFunc(t.f)
	handler.ServeHTTP(t.rr, t.req)
	return t
}
