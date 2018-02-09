package iw

import (
	"go-awt/infra/osencap"
	"go-awt/utils/iu"
	"fmt"
)

type NetworkCell struct {

}

var Iw = func (args ...string) (string, error)  {
	return osencap.Exec("iw", args...)
}

func Scan(iface string) ([]NetworkCell, error) {
	iface, _ = iu.ResolveIface(iface)
	out, err := Iw("dev", iface, "scan")
	if err != nil {
		return nil, err
	}
	fmt.Println(out)
	return nil, nil
}