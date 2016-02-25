package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/dashboard/healthcheck:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/dashboard/healthcheck:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

}
