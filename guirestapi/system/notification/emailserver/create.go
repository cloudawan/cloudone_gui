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

package emailserver

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type CreateController struct {
	beego.Controller
}

// @Title create
// @Description create the email server configuration
// @Param body body guirestapi.system.notification.emailserver.EmailServerSMTP true "body for email server configuration"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [post]
func (c *CreateController) Post() {
	inputBody := c.Ctx.Input.RequestBody
	emailServerSMTP := EmailServerSMTP{}
	err := json.Unmarshal(inputBody, &emailServerSMTP)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	}

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/emailserversmtp/"

	_, err = restclient.RequestPostWithStructure(url, emailServerSMTP, nil)

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJson()
		return
	} else {
		c.Data["json"] = make(map[string]interface{})
		c.ServeJson()
	}
}
