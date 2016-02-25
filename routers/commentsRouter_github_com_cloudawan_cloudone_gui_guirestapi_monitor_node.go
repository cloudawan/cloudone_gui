package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/node:DataController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/node:DataController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/node:IndexController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/node:IndexController"],
		beego.ControllerComments{
			"Get",
			`/information`,
			[]string{"get"},
			nil})

}
