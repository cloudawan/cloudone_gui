package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/system/namespace:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/system/namespace:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/system/namespace:EditController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/system/namespace:EditController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/system/namespace:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/system/namespace:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/system/namespace:SelectController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/system/namespace:SelectController"],
		beego.ControllerComments{
			"Put",
			`/select/:name`,
			[]string{"put"},
			nil})

}
