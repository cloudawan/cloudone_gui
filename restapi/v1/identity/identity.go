package identity

import (
	"encoding/json"
	"github.com/astaxie/beego/context"
)

const (
	loginURL = "/api/v1/identity/login"
)

func FilterToken(ctx *context.Context) {
	if (ctx.Input.IsGet() || ctx.Input.IsPost()) && ctx.Input.Url() == loginURL {
		// Don't redirect itself to prevent the circle
	} else {
		token := ctx.Input.Header("token")
		err := verifyToken(token)
		if err != nil {
			jsonMap := make(map[string]interface{})
			jsonMap["error"] = err.Error()
			byteSlice, _ := json.Marshal(jsonMap)
			ctx.Output.SetStatus(401)
			ctx.Output.Body(byteSlice)
			//ctx.Redirect(302, loginURL)
		}
	}
}
