package identity

import (
	"github.com/astaxie/beego/context"
)

const (
	loginPageURL = "/gui/login"
)

func FilterUser(ctx *context.Context) {
	if (ctx.Input.IsGet() || ctx.Input.IsPost()) && ctx.Input.Url() == loginPageURL {
		// Don't redirect itself to prevent the circle
	} else {
		_, ok := ctx.Input.Session("username").(string)
		if ok == false {
			ctx.Redirect(302, loginPageURL)
		}
	}
}
