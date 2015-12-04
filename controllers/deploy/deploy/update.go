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
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"strings"
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

func (c *UpdateController) Get() {
	c.TplNames = "deploy/deploy/update.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")
	oldVersion := c.GetString("oldVersion")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imagerecords/" + name

	imageRecordSlice := make([]ImageRecord, 0)

	_, err := restclient.RequestGetWithStructure(url, &imageRecordSlice)
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		versionEnvironmentMap := make(map[string]map[string]string)
		versionSlice := make([]string, 0)
		for _, imageRecord := range imageRecordSlice {
			if imageRecord.Version != oldVersion {
				versionSlice = append(versionSlice, imageRecord.Version)
				versionEnvironmentMap[imageRecord.Version] = imageRecord.Environment
			}
		}

		c.Data["name"] = name
		c.Data["versionSlice"] = versionSlice
		c.Data["versionEnvironmentMap"] = versionEnvironmentMap
	}

	guimessage.OutputMessage(c.Data)
}

func (c *UpdateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	namespaces, _ := c.GetSession("namespace").(string)

	imageInformationName := c.GetString("name")
	version := c.GetString("version")
	description := c.GetString("description")

	keySlice := make([]string, 0)
	inputMap := c.Input()
	if inputMap != nil {
		for key, _ := range inputMap {
			// Only collect environment belonging to this version
			if strings.HasPrefix(key, version) {
				keySlice = append(keySlice, key)
			}
		}
	}

	environmentSlice := make([]ReplicationControllerContainerEnvironment, 0)
	length := len(version) + 1 // + 1 for _
	for _, key := range keySlice {
		value := c.GetString(key)
		if len(value) > 0 {
			environmentSlice = append(environmentSlice,
				ReplicationControllerContainerEnvironment{key[length:], value})
		}
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/update/" + namespaces + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	deployUpdateInput := DeployUpdateInput{imageInformationName, version, description, environmentSlice}

	_, err := restclient.RequestPutWithStructure(url, deployUpdateInput, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Update deploy " + imageInformationName + " to version " + version + " success")
	}

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/deploy/deploy/")

	guimessage.RedirectMessage(c)
}
