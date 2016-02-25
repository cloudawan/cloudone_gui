package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/historicalcontainer:DataController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/historicalcontainer:DataController"],
		beego.ControllerComments{
			"Get",
			`/:replicationController`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/historicalcontainer:IndexController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/monitor/historicalcontainer:IndexController"],
		beego.ControllerComments{
			"Get",
			`/selectinformation`,
			[]string{"get"},
			nil})

}
