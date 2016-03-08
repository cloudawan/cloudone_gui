// Copyright 2015 CloudAwan LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	c.ServeJSON()
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
		c.ServeJSON()
		c.Abort("401")
		return
	}

	tokenString, err := createToken(loginRequestInput.Username, loginRequestInput.Password)

	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.ServeJSON()
		c.Abort("401")
		return
	}

	jsonMap := make(map[string]interface{})
	jsonMap["token"] = tokenString
	c.Data["json"] = jsonMap

	c.ServeJSON()
}
