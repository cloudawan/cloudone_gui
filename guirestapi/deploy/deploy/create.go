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
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type CreateController struct {
	beego.Controller
}

type DeployCreateInput struct {
	ImageInformationName string
	Version              string
	Description          string
	ReplicaAmount        int
	PortSlice            []ReplicationControllerContainerPort
	EnvironmentSlice     []ReplicationControllerContainerEnvironment
}

type ReplicationControllerContainerPort struct {
	Name          string
	ContainerPort int
}

type ReplicationControllerContainerEnvironment struct {
	Name  string
	Value string
}

type ImageRecord struct {
	ImageInformation string
	Version          string
	Path             string
	VersionInfo      map[string]string
	Environment      map[string]string
	Description      string
	CreatedTime      string
}

// @Title get
// @Description get the related selection in order to create a new deployment
// @Param name path string true "The name of image record"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router /createinformation/:name [get]
func (c *CreateController) Get() {
	name := c.GetString(":name")

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
			versionSlice = append(versionSlice, imageRecord.Version)
			versionEnvironmentMap[imageRecord.Version] = imageRecord.Environment
		}

		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["imageInformationName"] = name
		c.Data["json"].(map[string]interface{})["versionSlice"] = versionSlice
		c.Data["json"].(map[string]interface{})["versionEnvironmentMap"] = versionEnvironmentMap
		c.ServeJSON()
	}
}

// @Title create
// @Description create the autoscaler
// @Param body body guirestapi.deploy.deploy.DeployCreateInput true "body for deploy"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [post]
func (c *CreateController) Post() {
	inputBody := c.Ctx.Input.RequestBody
	deployCreateInput := DeployCreateInput{}
	err := json.Unmarshal(inputBody, &deployCreateInput)
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

	namespaces, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/create/" + namespaces + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPostWithStructure(url, deployCreateInput, nil, tokenHeaderMap)

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
