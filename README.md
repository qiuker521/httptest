# httptest

[中文文档](README_zh.md)


[![GoDoc](https://godoc.org/github.com/qiuker521/httptest?status.svg)](https://godoc.org/github.com/qiuker521/httptest)
[![Build Status](https://travis-ci.org/qiuker521/httptest.svg?branch=master)](https://travis-ci.org/qiuker521/httptest)
[![Go Report Card](https://goreportcard.com/badge/github.com/qiuker521/httptest)](https://goreportcard.com/report/github.com/qiuker521/httptest)


### What we want ?

In development we found the advantages of unit tests, 
and golang it self provides good test tools.

However, some test tools are too board, we have to add a layer for using it.

As well as the [gorequest](https://github.com/parnurzeal/gorequest) does to `http.Request`, 
we add a layer on httptest to make http unit test easy.


### How to use it?

As the [Testing Your (HTTP) Handlers in Go](https://elithrar.github.io/article/testing-http-handlers-go/) says, write a http unit for a http handler is a boring work.

We just simplify it into some several lines:


```go
//badHandler is a handler always returns 400
func badHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not a regular name or password", http.StatusBadRequest)
}

//We want to test if it returns 400, we could write:
New("/bad", badHandler, t).Do().CheckCode(http.StatusBadRequest)

//So the 200 is as follows:(this is always fail)
New("/ok", badHandler, t).Do().CheckCode(http.StatusOK)

//Add header to request, and just do a request, using POST:
New("/", badHandler, t).Post().AddParams("name", "value1").AddParams("nam22", "value3").Do()

//Add cookie to request, and test if the response body contains a certain string:
New("/", cookieHandler, t).Get().AddCookies(cookie).Do().BodyContains("testcookievalue")

//Just get the *http.ResponseRecorder, do every thing your self.
rr = New("/dump", headerHandler, t).Post().AddParams("name", "value1").Do().ResponseRecorder()

//We forget to add parameter to request
New("/ok", badHandler, t).AddParams("a", "aa").AddParams("b", "bb").Do().CheckCode(http.StatusOK)

//Use http basic auth:
New("/bad", badHandler, t).SetBasicAuth(username, password).Do().CheckCode(http.StatusBadRequest)

//Use your self-defined http.Request:
New("/bad", badHandler, t).SetRequest(req).Do().CheckCode(http.StatusBadRequest)

//And more in test file and source code.

```

Don't forget `.Do()`.


### What we did and roadmap?

- [x] Add common layer for request.
- [x] Add common header, parameter for request.
- [x] Do wrapper on httptest.ResponseRecorder
- [ ] Do better wrapper on body test, such as json.
