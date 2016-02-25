package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploybluegreen:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploybluegreen:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:imageInformation`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploybluegreen:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploybluegreen:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploybluegreen:SelectController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploybluegreen:SelectController"],
		beego.ControllerComments{
			"Get",
			`/selectinformation/:imageInformation`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploybluegreen:SelectController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deploybluegreen:SelectController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

}
