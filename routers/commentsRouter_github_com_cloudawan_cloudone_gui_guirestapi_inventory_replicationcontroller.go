package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:namespace/:replicationcontroller`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:EditController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:EditController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:PodLogController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:PodLogController"],
		beego.ControllerComments{
			"Get",
			`/podlog/:namespace/:pod`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:SizeController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/inventory/replicationcontroller:SizeController"],
		beego.ControllerComments{
			"Put",
			`/resize/:name`,
			[]string{"put"},
			nil})

}
