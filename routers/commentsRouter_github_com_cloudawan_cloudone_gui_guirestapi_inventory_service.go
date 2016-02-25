package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/service:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/service:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:namespace/:service`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/service:EditController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/service:EditController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/service:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/service:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

}
