package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/autoscaler:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/autoscaler:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:kind/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/autoscaler:EditController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/autoscaler:EditController"],
		beego.ControllerComments{
			"Get",
			`/:kind/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/autoscaler:EditController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/autoscaler:EditController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/autoscaler:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/deploy/autoscaler:ListController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

}
