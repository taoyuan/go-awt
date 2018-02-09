package iw

import (
	"go-awt/infra/osencap"
	"go-awt/utils/iu"
	"fmt"
)

type NetworkCell struct {

}

func Scan(iface string) ([]NetworkCell, error) {
	if iface == "" {
		iface = "default"
	}
	iface, _ = iu.ResolveIface(iface)
	args := []string{"dev", iface, "scan"}
	out, err := osencap.Exec("iw", args...)
	if err != nil {
		return nil, err
	}
	fmt.Println(out)
	return nil, nil
}