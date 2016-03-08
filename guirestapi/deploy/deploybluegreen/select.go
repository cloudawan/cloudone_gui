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

package deploybluegreen

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type SelectController struct {
	beego.Controller
}

// @Title get
// @Description get the related selection in order to select a blue green deployment
// @Param imageInformation path string true "The name of image record"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router /selectinformation/:imageInformation [get]
func (c *SelectController) Get() {
	imageInformation := c.GetString(":imageInformation")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploybluegreens/deployable/" + imageInformation + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	namespaceSlice := make([]string, 0)
	_, err := restclient.RequestGetWithStructure(url, &namespaceSlice)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["namespaceSlice"] = namespaceSlice
		c.ServeJSON()
	}
}

// @Title select
// @Description select the different namespace
// @Param body body guirestapi.deploy.deploybluegreen.DeployBlueGreen true "body for blue green deployment"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [put]
func (c *SelectController) Put() {
	inputBody := c.Ctx.Input.RequestBody
	deployBlueGreen := DeployBlueGreen{}
	err := json.Unmarshal(inputBody, &deployBlueGreen)
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
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploybluegreens/" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	_, err = restclient.RequestPutWithStructure(url, deployBlueGreen, nil)

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
