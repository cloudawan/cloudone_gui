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

package replicationcontroller

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/limit"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type EditController struct {
	beego.Controller
}

type ReplicationController struct {
	Name           string
	ReplicaAmount  int
	Selector       ReplicationControllerSelector
	Label          ReplicationControllerLabel
	ContainerSlice []ReplicationControllerContainer
}

type ReplicationControllerSelector struct {
	Name    string
	Version string
}

type ReplicationControllerLabel struct {
	Name string
}

type ReplicationControllerContainer struct {
	Name      string
	Image     string
	PortSlice []ReplicationControllerContainerPort
}

type ReplicationControllerContainerPort struct {
	Name          string
	ContainerPort int
}

// @Title create
// @Description create replication controller
// @Param body body guirestapi.inventory.replicationcontroller.ReplicationController true "body for replication controller"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [post]
func (c *EditController) Post() {
	inputBody := c.Ctx.Input.CopyBody(limit.InputPostBodyMaximum)
	replicationController := ReplicationController{}
	err := json.Unmarshal(inputBody, &replicationController)
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

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/replicationcontrollers/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPostWithStructure(url, replicationController, nil, tokenHeaderMap)

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
