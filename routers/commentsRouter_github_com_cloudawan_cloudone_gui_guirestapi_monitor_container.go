package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/container:DataController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/container:DataController"],
		beego.ControllerComments{
			"Get",
			`/:replicationController`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/container:IndexController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/container:IndexController"],
		beego.ControllerComments{
			"Get",
			`/selectinformation`,
			[]string{"get"},
			nil})

}
