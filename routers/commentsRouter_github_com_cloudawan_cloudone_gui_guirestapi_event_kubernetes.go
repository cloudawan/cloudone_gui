package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/event/kubernetes:AcknowledgeController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/event/kubernetes:AcknowledgeController"],
		beego.ControllerComments{
			"Put",
			`/acknowledge/:namespace/:id`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/event/kubernetes:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/event/kubernetes:ListController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

}
