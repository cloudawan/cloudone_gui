package identity

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

type LoginRequestInput struct {
	Username string
	Password string
}

func (c *LoginController) Get() {
	errorJsonMap := make(map[string]interface{})
	errorJsonMap["error"] = "Unauthorized"
	c.Data["json"] = errorJsonMap
	c.ServeJson()
	c.Abort("401")
}

func (c *LoginController) Post() {
	loginRequestInput := LoginRequestInput{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &loginRequestInput)
	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.ServeJson()
		c.Abort("401")
		return
	}

	tokenString, err := createToken(loginRequestInput.Username, loginRequestInput.Password)

	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.ServeJson()
		c.Abort("401")
		return
	}

	jsonMap := make(map[string]interface{})
	jsonMap["token"] = tokenString
	c.Data["json"] = jsonMap

	c.ServeJson()
}
