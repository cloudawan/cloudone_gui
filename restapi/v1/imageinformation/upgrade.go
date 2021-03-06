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

package imageinformation

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/limit"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type UpdateController struct {
	beego.Controller
}

type DeployUpgradeInput struct {
	ImageInformationName string
	Description          string
}

func (c *UpdateController) Post() {
	inputBody := c.Ctx.Input.CopyBody(limit.InputPostBodyMaximum)

	token := c.Ctx.Input.Header("Token")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	deployUpgradeInput := DeployUpgradeInput{}
	err := json.Unmarshal(inputBody, &deployUpgradeInput)
	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.Ctx.Output.Status = 400
		c.ServeJSON()
		return
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imageinformations/upgrade/"

	tokenHeaderMap := make(map[string]string, 0)
	tokenHeaderMap["token"] = token

	_, err = restclient.RequestPutWithStructure(url, deployUpgradeInput, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.Ctx.Output.Status = 400
		c.ServeJSON()
		return
	}

	jsonMap := make(map[string]interface{})
	c.Data["json"] = jsonMap
	c.ServeJSON()
}
