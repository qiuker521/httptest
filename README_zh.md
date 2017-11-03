# httptest


### What we want ?

不扯淡，简化httptest库。


### 用法

标准httptest库用法可以参考[Testing Your (HTTP) Handlers in Go](https://elithrar.github.io/article/testing-http-handlers-go/)，但是有目共睹，写一个单元测试特别长，所以我们简化了一下。


```go
//一个永远返回400的测试handler
func badHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not a regular name or password", http.StatusBadRequest)
}

//测试这个handler是否返回400
New("/bad", badHandler, t).Do().CheckCode(http.StatusBadRequest)

//测试他是不是返回200（当然会测试失败）
New("/ok", badHandler, t).Do().CheckCode(http.StatusOK)

//带着header测试
New("/", badHandler, t).Post().AddParams("name", "value1").AddParams("nam22", "value3").Do()

//带着cookie测试，并且判断结果是否包含字符串。
New("/", cookieHandler, t).Get().AddCookies(cookie).Do().BodyContains("testcookievalue")

//获取 *http.ResponseRecorder, 然后自己测试
rr = New("/dump", headerHandler, t).Post().AddParams("name", "value1").Do().ResponseRecorder()

//给请求加参数，不写默认是GET请求
New("/ok", badHandler, t).AddParams("a", "aa").AddParams("b", "bb").Do().CheckCode(http.StatusOK)

//http basic auth:
New("/bad", badHandler, t).SetBasicAuth(username, password).Do().CheckCode(http.StatusBadRequest)

//自己定制 http.Request:
New("/bad", badHandler, t).SetRequest(req).Do().CheckCode(http.StatusBadRequest)

//And more in test file and source code.

```

必须有 `.Do()`，才能进行请求，不然不会请求。

### 后续

[] 另外，还有计划对获取到的请求进行json解析……
