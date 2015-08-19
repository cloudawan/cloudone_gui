package controllers

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	c.Ctx.Redirect(302, "/gui/login")

	guimessage.RedirectMessage(c)
}
