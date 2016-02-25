package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:EditController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:EditController"],
		beego.ControllerComments{
			"Get",
			`/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:EditController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:EditController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:LaunchController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:LaunchController"],
		beego.ControllerComments{
			"Get",
			`/launchinformation/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:LaunchController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:LaunchController"],
		beego.ControllerComments{
			"Post",
			`/launch/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/thirdparty:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

}
