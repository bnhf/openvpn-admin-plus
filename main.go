package main

import (
	"github.com/astaxie/beego"
	"github.com/bnhf/openvpn-tap-external-web-ui/lib"
	_ "github.com/bnhf/openvpn-tap-external-web-ui/routers"
)

func main() {
	lib.AddFuncMaps()
	beego.Run()
}
