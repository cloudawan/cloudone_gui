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

package autoscaler

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type EditController struct {
	beego.Controller
}

// @Title get
// @Description get the autoscaler
// @Param kind path string true "The type of target autoscaler configured for"
// @Param name path string true "The name of target autoscaler configured for"
// @Success 200 {object} guirestapi.deploy.autoscaler.ReplicationControllerAutoScaler
// @Failure 404 error reason
// @router /:kind/:name [get]
func (c *EditController) Get() {
	kind := c.GetString(":kind")
	name := c.GetString(":name")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespace, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/autoscalers/" + namespace + "/" + kind + "/" + name

	replicationControllerAutoScaler := ReplicationControllerAutoScaler{}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &replicationControllerAutoScaler, tokenHeaderMap)

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
		c.Data["json"] = replicationControllerAutoScaler
		c.ServeJSON()
	}
}

// @Title update
// @Description update the autoscaler
// @Param body body guirestapi.deploy.autoscaler.ReplicationControllerAutoScaler true "body for autoscaler"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [put]
func (c *EditController) Put() {
	inputBody := c.Ctx.Input.RequestBody
	replicationControllerAutoScaler := ReplicationControllerAutoScaler{}
	err := json.Unmarshal(inputBody, &replicationControllerAutoScaler)
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

	namespace, _ := c.GetSession("namespace").(string)

	replicationControllerAutoScaler.KubeapiHost = kubeapiHost
	replicationControllerAutoScaler.KubeapiPort = kubeapiPort
	replicationControllerAutoScaler.Namespace = namespace

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/autoscalers/"

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPutWithStructure(url, replicationControllerAutoScaler, nil, tokenHeaderMap)

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
