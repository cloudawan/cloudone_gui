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

package deployclusterapplication

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/limit"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type Cluster struct {
	Name                      string
	Description               string
	ReplicationControllerJson string
	ServiceJson               string
	Environment               map[string]string
	ScriptType                string
	ScriptContent             string
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
	Name             string
	Image            string
	PortSlice        []ReplicationControllerContainerPort
	EnvironmentSlice []ReplicationControllerContainerEnvironment
}

type ReplicationControllerContainerPort struct {
	Name          string
	ContainerPort int
}

type ReplicationControllerContainerEnvironment struct {
	Name  string
	Value string
}

type Environment struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SizeController struct {
	beego.Controller
}

// @Title get
// @Description get the related selection in order to change size
// @Param name path string true "The name of cluster application deployment"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router /sizeinformation/:name [get]
func (c *SizeController) Get() {
	name := c.GetString(":name")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	namespace, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/clusterapplications/" + name
	cluster := Cluster{}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &cluster, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = "Fail to get cluster application with error" + err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}

	c.Data["name"] = name

	// Get configured environment from any one of replication controller belonging to this cluster application
	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deployclusterapplications/" + namespace + "/" + name
	deployClusterApplication := DeployClusterApplication{}

	_, err = restclient.RequestGetWithStructure(url, &deployClusterApplication, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = "Fail to get cluster application deployment with error" + err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}

	if len(deployClusterApplication.ReplicationControllerNameSlice) == 0 {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = "The replication controller name slice is empty for the cluster application deployment with name " + name
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}

	replicationControllerName := deployClusterApplication.ReplicationControllerNameSlice[0]

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/replicationcontrollers/" + namespace + "/" + replicationControllerName +
		"?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)
	replicationController := ReplicationController{}

	_, err = restclient.RequestGetWithStructure(url, &replicationController, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = "Fail to get the replication controller with name " + replicationControllerName
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}

	for _, container := range replicationController.ContainerSlice {
		for _, environment := range container.EnvironmentSlice {
			cluster.Environment[environment.Name] = environment.Value
		}
	}

	c.Data["json"] = make(map[string]interface{})
	c.Data["json"].(map[string]interface{})["environment"] = cluster.Environment
	c.ServeJSON()
}

// @Title resize
// @Description resize the cluster application
// @Param body body []Environment true "Array of environment"
// @Param name path string true "The name of cluster application deployment"
// @Param size query string true "The size to change"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router /size/:name [put]
func (c *SizeController) Put() {
	inputBody := c.Ctx.Input.CopyBody(limit.InputPostBodyMaximum)
	environmentSlice := make([]Environment, 0)
	err := json.Unmarshal(inputBody, &environmentSlice)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}

	name := c.GetString(":name")
	size := c.GetString("size")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	namespace, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deployclusterapplications/size/" + namespace + "/" + name +
		"?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort) + "&size=" + size

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPut(url, environmentSlice, tokenHeaderMap, true)

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
