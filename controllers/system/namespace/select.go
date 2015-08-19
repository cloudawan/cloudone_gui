package namespace

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
)

type SelectController struct {
	beego.Controller
}

func (c *SelectController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	name := c.GetString("name")

	c.SetSession("namespace", name)
	guimessage.AddSuccess("Use namespace " + name)

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/system/namespace/")

	guimessage.RedirectMessage(c)
}
