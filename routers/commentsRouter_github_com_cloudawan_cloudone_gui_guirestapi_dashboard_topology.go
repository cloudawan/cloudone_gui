package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/dashboard/topology:IndexController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/dashboard/topology:IndexController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

}
