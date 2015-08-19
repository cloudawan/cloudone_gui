package routers

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/restapi/v1/identity"
	"github.com/cloudawan/kubernetes_management_gui/restapi/v1/imageinformation"
)

func init() {
	beego.Router("/api/v1/identity/login", &identity.LoginController{})
	beego.Router("/api/v1/imageinformation/upgrade", &imageinformation.UpdateController{})
}
