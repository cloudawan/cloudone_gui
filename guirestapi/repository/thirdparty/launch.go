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

package thirdparty

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/limit"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strings"
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

type ClusterLaunch struct {
	Size                              int
	EnvironmentSlice                  []interface{}
	ReplicationControllerExtraJsonMap map[string]interface{}
}

type LaunchController struct {
	beego.Controller
}

// @Title get
// @Description get the cluster application template
// @Param name path string true "The name of cluster application template"
// @Success 200 {object} guirestapi.repository.thirdparty.Cluster
// @Failure 404 error reason
// @router /launchinformation/:name [get]
func (c *LaunchController) Get() {
	name := c.GetString("name")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

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
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = cluster
		c.ServeJSON()
	}
}

// @Title launch
// @Description launch the cluster application template
// @Param body body guirestapi.repository.thirdparty.Cluster true "body for cluster application template"
// @Param name query string true "The name to use"
// @Param size query string true "The size to use"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router /launch/ [post]
func (c *LaunchController) Post() {
	name := c.GetString("name")
	size, _ := c.GetInt("size")

	inputBody := c.Ctx.Input.CopyBody(limit.InputPostBodyMaximum)
	environmentSlice := make([]interface{}, 0)
	err := json.Unmarshal(inputBody, &environmentSlice)
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

	namespace, _ := c.GetSession("namespace").(string)

	clusterLaunch := ClusterLaunch{
		size,
		environmentSlice,
		nil,
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/clusterapplications/launch/" + namespace + "/" + name
	jsonMap := make(map[string]interface{})

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPostWithStructure(url, clusterLaunch, &jsonMap, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		errorMessage, _ := jsonMap["Error"].(string)
		if strings.HasPrefix(errorMessage, "Replication controller already exists") {
			// Error
			c.Data["json"] = make(map[string]interface{})
			c.Data["json"].(map[string]interface{})["error"] = "Replication controller " + name + " already exists"
			c.Ctx.Output.Status = 404
			c.ServeJSON()
			return
		} else {
			// Error
			c.Data["json"] = make(map[string]interface{})
			c.Data["json"].(map[string]interface{})["error"] = err.Error()
			c.Ctx.Output.Status = 404
			c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = make(map[string]interface{})
		c.ServeJSON()
	}
}
