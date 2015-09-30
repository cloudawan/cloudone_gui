// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     Unknwon
package main

import (
	"github.com/cloudawan/kubernetes_management_gui/Godeps/_workspace/src/github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/Godeps/_workspace/src/github.com/astaxie/beego/example/chat/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WSController{})
	beego.Run()
}
