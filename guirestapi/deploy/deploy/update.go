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

package deploy

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

type DeployUpdateInput struct {
	ImageInformationName string
	Version              string
	Description          string
	EnvironmentSlice     []ReplicationControllerContainerEnvironment
}

// @Title get
// @Description get the related selection in order to create a new deployment
// @Param name path string true "The name of image record"
// @Param oldVersion query string true "The current version to be replaced"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router /updateinformation/:name [get]
func (c *UpdateController) Get() {
	name := c.GetString(":name")
	oldVersion := c.GetString("oldVersion")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imagerecords/" + name

	imageRecordSlice := make([]ImageRecord, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &imageRecordSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	} else {
		versionEnvironmentMap := make(map[string]map[string]string)
		versionSlice := make([]string, 0)
		for _, imageRecord := range imageRecordSlice {
			if imageRecord.Version != oldVersion {
				versionSlice = append(versionSlice, imageRecord.Version)
				versionEnvironmentMap[imageRecord.Version] = imageRecord.Environment
			}
		}

		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["name"] = name
		c.Data["json"].(map[string]interface{})["versionSlice"] = versionSlice
		c.Data["json"].(map[string]interface{})["versionEnvironmentMap"] = versionEnvironmentMap
		c.ServeJSON()
	}
}

// @Title update
// @Description update the deploy
// @Param body body guirestapi.deploy.deploy.DeployUpdateInput true "body for deploy"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [put]
func (c *UpdateController) Put() {
	inputBody := c.Ctx.Input.CopyBody(limit.InputPostBodyMaximum)
	deployUpdateInput := DeployUpdateInput{}
	err := json.Unmarshal(inputBody, &deployUpdateInput)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespaces, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/update/" + namespaces

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPutWithStructure(url, deployUpdateInput, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = make(map[string]interface{})
		c.ServeJSON()
	}
}
