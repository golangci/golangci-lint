//golangcitest:args -Enoctx
package testdata

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import (
	"context"
	"net/http"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

var newRequestPkg = http.NewRequest

func _() {
	const url = "http://example.com"
	cli := &http.Client{}

	ctx := context.Background()
	http.Get(url) // want `net/http\.Get must not be called`
	_ = http.Get  // OK
	f := http.Get // OK
	f(url)        // want `net/http\.Get must not be called`

	http.Head(url)          // want `net/http\.Head must not be called`
	http.Post(url, "", nil) // want `net/http\.Post must not be called`
	http.PostForm(url, nil) // want `net/http\.PostForm must not be called`

	cli.Get(url) // want `\(\*net/http\.Client\)\.Get must not be called`
	_ = cli.Get  // OK
	m := cli.Get // OK
	m(url)       // want `\(\*net/http\.Client\)\.Get must not be called`

	cli.Head(url)          // want `\(\*net/http\.Client\)\.Head must not be called`
	cli.Post(url, "", nil) // want `\(\*net/http\.Client\)\.Post must not be called`
	cli.PostForm(url, nil) // want `\(\*net/http\.Client\)\.PostForm must not be called`

	req, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	cli.Do(req)

	req2, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil) // OK
	cli.Do(req2)

	req3, _ := http.NewRequest(http.MethodPost, url, nil) // OK
	req3 = req3.WithContext(ctx)
	cli.Do(req3)

	f2 := func(req *http.Request, ctx context.Context) *http.Request {
		return req
	}
	req4, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	req4 = f2(req4, ctx)

	req41, _ := http.NewRequest(http.MethodPost, url, nil) // OK
	req41 = req41.WithContext(ctx)
	req41 = f2(req41, ctx)

	newRequest := http.NewRequest
	req5, _ := newRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	cli.Do(req5)

	req51, _ := newRequest(http.MethodPost, url, nil) // OK
	req51 = req51.WithContext(ctx)
	cli.Do(req51)

	req52, _ := newRequestPkg(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	cli.Do(req52)

	type MyRequest = http.Request
	f3 := func(req *MyRequest, ctx context.Context) *MyRequest {
		return req
	}
	req6, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	req6 = f3(req6, ctx)

	req61, _ := http.NewRequest(http.MethodPost, url, nil) // OK
	req61 = req61.WithContext(ctx)
	req61 = f3(req61, ctx)

	type MyRequest2 http.Request
	f4 := func(req *MyRequest2, ctx context.Context) *MyRequest2 {
		return req
	}
	req7, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	req71 := MyRequest2(*req7)
	f4(&req71, ctx)

	req72, _ := http.NewRequest(http.MethodPost, url, nil) // OK
	req72 = req72.WithContext(ctx)
	req73 := MyRequest2(*req7)
	f4(&req73, ctx)

	req8, _ := func() (*http.Request, error) {
		return http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	}()
	cli.Do(req8)

	req82, _ := func() (*http.Request, error) {
		req82, _ := http.NewRequest(http.MethodPost, url, nil) // OK
		req82 = req82.WithContext(ctx)
		return req82, nil
	}()
	cli.Do(req82)

	f5 := func(req, req2 *http.Request, ctx context.Context) (*http.Request, *http.Request) {
		return req, req2
	}
	req9, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	req9, _ = f5(req9, req9, ctx)

	req91, _ := http.NewRequest(http.MethodPost, url, nil) // OK
	req91 = req91.WithContext(ctx)
	req9, _ = f5(req91, req91, ctx)

	req10, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	req11, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	req10, req11 = f5(req10, req11, ctx)

	req101, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
	req111, _ := http.NewRequest(http.MethodPost, url, nil) // OK
	req111 = req111.WithContext(ctx)
	req101, req111 = f5(req101, req111, ctx)

	func() (*http.Request, *http.Request) {
		req12, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
		req13, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
		return req12, req13
	}()

	func() (*http.Request, *http.Request) {
		req14, _ := http.NewRequest(http.MethodPost, url, nil) // want `should rewrite http.NewRequestWithContext or add \(\*Request\).WithContext`
		req15, _ := http.NewRequest(http.MethodPost, url, nil) // OK
		req15 = req15.WithContext(ctx)

		return req14, req15
	}()
}
