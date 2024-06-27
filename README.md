# Supported
- [x] HTTP
- [x] HTTPS
- [x] SOCKS4
- [x] SOCKS4A
- [x] SOCKS5
- [x] SSH
- [x] ShadowSocks
- [x] VMess

# Example
```go
package main

import (
	"fmt"
	"github.com/dneht/proxyclient"
	"io"
	"net/http"
)

func main() {
	//http://localhost:1800
	//socks://localhost:1800
	//ss://localhost:1800
	//vmess://localhost:1800
	proxy, err := proxyclient.New("ssh://localhost:22", &proxyclient.Option{
		User: &proxyclient.UserOption{},
		SSH: &proxyclient.SSHOption{
			Name: "user",
		},
		SS: &proxyclient.SSOption{
			Method:   "cipher",
			Password: "pwd",
		},
		VMess: &proxyclient.VMessOption{
			UUID:     "uuid",
			Security: "auto",
			AlterId:  0,
		},
	})
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: proxy.DialContext,
		},
	}
	request, err := client.Get("https://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	content, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(content))
}
```