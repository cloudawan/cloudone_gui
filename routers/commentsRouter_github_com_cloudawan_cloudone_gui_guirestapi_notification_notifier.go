package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/notification/notifier:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/notification/notifier:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:kind/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/notification/notifier:EditController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/notification/notifier:EditController"],
		beego.ControllerComments{
			"Get",
			`/:kind/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/notification/notifier:EditController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/notification/notifier:EditController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/notification/notifier:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/notification/notifier:ListController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

}
