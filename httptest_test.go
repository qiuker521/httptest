package httptest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func badHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not a regular name or password", http.StatusBadRequest)
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not a regular name or password", http.StatusOK)
}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)
	for name := range r.Form {
		w.Header().Set(name, r.FormValue(name))
	}
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	content := ""
	for name := range r.Form {
		content += name
		content += r.FormValue(name)
	}
	w.Write([]byte(content))
}

func TestCheckCode(t *testing.T) {
	New("/bad", badHandler, t).Do().CheckCode(http.StatusBadRequest)
	New("/ok", okHandler, t).Do().CheckCode(http.StatusOK)
}

func TestCheckHeader(t *testing.T) {
	New("/dump", headerHandler, t).
		Post().
		AddParams("name", "value1").
		AddParams("nam22", "value3").
		Do().
		CheckCode(http.StatusOK).
		CheckHeader("name", "value1").
		CheckHeader("nam22", "value3")
}

func TestBodyContains(t *testing.T) {
	New("/content", contentHandler, t).
		Post().
		AddParams("name", "value1").
		AddParams("nam22", "value3").
		Do().
		CheckCode(http.StatusOK).
		CheckBodyContains("nam22").
		CheckBodyContains("value3")
}

func TestResponseRecorder(t *testing.T) {
	var rr *httptest.ResponseRecorder
	rr = New("/content", contentHandler, t).
		Get().
		AddParams("name", "value1").
		AddParams("nam22", "value3").
		Do().
		GetResponseRecorder()
	if rr.Code != http.StatusOK {
		t.Errorf("error while get response code: want [%d], got [%d]", http.StatusOK, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "nam22") {
		t.Errorf("error while get body, want [%s], got none", "nam22")
	}
	rr = New("/dump", headerHandler, t).
		Get().
		AddParams("name", "value1").
		AddParams("nam22", "value3").
		Do().
		GetResponseRecorder()

	if rr.Header().Get("nam22") != "value3" {
		t.Errorf("want header [%s] to equal: [%s], but got: [%s]", "nam22", "value3", rr.Header().Get("nam22"))
	}

	rr = New("/content", contentHandler, t).
		Post().
		AddParams("name", "value1").
		AddParams("nam22", "value3").
		Do().
		GetResponseRecorder()
	if rr.Code != http.StatusOK {
		t.Errorf("error while get response code: want [%d], got [%d]", http.StatusOK, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "nam22") {
		t.Errorf("error while get body, want [%s], got none", "nam22")
	}
	rr = New("/dump", headerHandler, t).
		Post().
		AddParams("name", "value1").
		AddParams("nam22", "value3").
		Do().
		GetResponseRecorder()

	if rr.Header().Get("nam22") != "value3" {
		t.Errorf("want header [%s] to equal: [%s], but got: [%s]", "nam22", "value3", rr.Header().Get("nam22"))
	}
}

func cookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("testcookiename")
	w.Write([]byte(cookie.Value))
}

func TestCookie(t *testing.T) {
	cookie := &http.Cookie{Name: "testcookiename", Value: "testcookievalue", Path: "/", MaxAge: 86400}
	New("/", cookieHandler, t).
		Get().
		AddCookies(cookie).
		Do().
		CheckBodyContains("testcookievalue")

	New("/", cookieHandler, t).
		Post().
		AddCookies(cookie).
		Do().
		CheckBodyContains("testcookievalue")
}
