package main

import (
	"github.com/taoyuan/go-awt/infra/osencap"
	"log"
)

func main() {
	out, err := osencap.Exec("bash", "utils/cap/test/create_ap_error")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(out)
}
