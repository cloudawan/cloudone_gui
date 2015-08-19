package identity

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
)

type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	c.DelSession("username")

	c.Ctx.Redirect(302, "/gui/login")

	guimessage.RedirectMessage(c)
}
