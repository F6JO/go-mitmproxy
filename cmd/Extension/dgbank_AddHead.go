package Extension

import "github.com/lqqyt2423/go-mitmproxy/proxy"

type AddHeader struct {
	proxy.BaseAddon
	count int
}

func (a *AddHeader) Requestheaders(f *proxy.Flow) {
	a.count += 1
	f.Request.Header.Add("Dgbank_Test_Scan", "True")
}
