// Package ipnetutil IP utils.
package ipnetutil

import (
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/xmopen/golib/pkg/xconfig"
)

var ipDataBuffer []byte

func init() {
	buffer, err := xdb.LoadContentFromFile(xconfig.ParseIPXDBConfigPath)
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
