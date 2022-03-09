package main

import (
	"github.com/astaxie/beego"
	"github.com/bnhf/pivpn-tap-web-ui/lib"
	_ "github.com/bnhf/pivpn-tap-web-ui/routers"
)

func main() {
	lib.AddFuncMaps()
	beego.Run()
}
