package identity

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"strconv"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplNames = "identity/login.html"
}

func (c *LoginController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	username := c.GetString("username")
	password := c.GetString("password")
	timeZoneOffset, err := c.GetInt("timeZoneOffset")
	if err != nil {
		guimessage.AddDanger("Fail to get browser time zone offset. Use UTC instead")
	} else {
		hourOffset := float64(timeZoneOffset) / 60.0
		sign := "-"
		if hourOffset < 0 {
			sign = "+"
		}
		guimessage.AddSuccess("Browser time zone is " + sign + strconv.FormatFloat(hourOffset, 'f', -1, 64) + " from UTC")
		c.SetSession("timeZoneOffset", timeZoneOffset)
	}

	// TODO RBAC
	if username == "admin" && password == "password" {
		// Username
		c.SetSession("username", username)
		// Namespace
		namespace := beego.AppConfig.String("namespace")
		c.SetSession("namespace", namespace)

		guimessage.AddSuccess("User " + username + " login")
	}

	c.Ctx.Redirect(302, "/gui/dashboard/topology/")

	guimessage.RedirectMessage(c)
}
