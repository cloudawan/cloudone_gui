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
	"github.com/cloudawan/cloudone_utility/restclient"
)

type LoginController struct {
	beego.Controller
}

type UserData struct {
	Username string
	Password string
}

type TokenData struct {
	Token string
}

func (c *LoginController) Get() {
	errorJsonMap := make(map[string]interface{})
	errorJsonMap["error"] = "Unauthorized"
	c.Data["json"] = errorJsonMap
	c.ServeJSON()
	c.Abort("401")
}

func (c *LoginController) Post() {
	userData := UserData{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userData)
	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["ErrorMessage"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.ServeJSON()
		c.Abort("401")
		return
	}

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	// User
	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/authorizations/tokens/"

	tokenData := TokenData{}
	errorInterface, err := restclient.RequestPostWithStructure(url, userData, &tokenData, nil)
	errorJsonMap, _ := errorInterface.(map[string]interface{})

	if err != nil {
		// Error
		c.Data["json"] = errorJsonMap
		c.ServeJSON()
		c.Abort("401")
		return
	}

	headerMap := make(map[string]string)
	headerMap["token"] = tokenData.Token
	c.SetSession("tokenHeaderMap", headerMap)

	jsonMap := make(map[string]interface{})
	jsonMap["Token"] = tokenData.Token
	c.Data["json"] = jsonMap

	c.ServeJSON()
}
