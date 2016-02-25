package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/imageinformation:CreateController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/imageinformation:CreateController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/imageinformation:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/imageinformation:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/imageinformation:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/imageinformation:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/imageinformation:UpgradeController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/repository/imageinformation:UpgradeController"],
		beego.ControllerComments{
			"Put",
			`/upgrade/`,
			[]string{"put"},
			nil})

}
