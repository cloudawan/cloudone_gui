package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deployclusterapplication:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deployclusterapplication:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:clusterApplicationName`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deployclusterapplication:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deployclusterapplication:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deployclusterapplication:SizeController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deployclusterapplication:SizeController"],
		beego.ControllerComments{
			"Get",
			`/sizeinformation/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deployclusterapplication:SizeController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/deployclusterapplication:SizeController"],
		beego.ControllerComments{
			"Put",
			`/size/:name`,
			[]string{"put"},
			nil})

}
