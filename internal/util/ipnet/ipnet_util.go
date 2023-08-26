// Package ipnetutil IP utils.
package ipnetutil

import (
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/xmopen/golib/pkg/xconfig"
)

// Linux要切换,后续切换到lib.
const xdbPath = xconfig.RuntimeENVWindowsGlobalConfigPath + "\\" + "ip2region.xdb"

var (
	ipDataBuffer []byte
)

func init() {
	buffer, err := xdb.LoadContentFromFile(xdbPath)
	if err != nil {
		panic(err)
	}
	ipDataBuffer = buffer
}

// ParseIPLocation parse ip location
func ParseIPLocation(ip string) (string, error) {
	search, err := xdb.NewWithBuffer(ipDataBuffer)
	if err != nil {
		return "", err
	}
	defer search.Close()
	region, err := search.SearchByStr(ip)
	if err != nil {
		return "", err
	}
	return region, nil
}
