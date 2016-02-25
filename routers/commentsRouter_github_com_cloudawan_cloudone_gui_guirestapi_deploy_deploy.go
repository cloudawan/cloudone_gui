package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:CreateController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:CreateController"],
		beego.ControllerComments{
			"Get",
			`/createinformation/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:CreateController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:CreateController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:UpdateController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:UpdateController"],
		beego.ControllerComments{
			"Get",
			`/updateinformation/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:UpdateController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploy:UpdateController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

}
