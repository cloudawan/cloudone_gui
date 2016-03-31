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

func (c *UpdateController) Get() {
	c.TplName = "deploy/deploy/update.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespace, _ := c.GetSession("namespace").(string)

	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		// Redirect to list
		c.Ctx.Redirect(302, "/gui/deploy/deploy/")
		return
	}

	name := c.GetString("name")
	oldVersion := c.GetString("oldVersion")

	replicationControllerName := name + oldVersion

	// Retrieve the current environment setting
	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/replicationcontrollers/" + namespace + "/" + replicationControllerName +
		"?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)
	replicationController := ReplicationController{}
	_, err = restclient.RequestGetWithStructure(url, &replicationController)
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		// Redirect to list
		c.Ctx.Redirect(302, "/gui/deploy/deploy/")
		return
	}

	// Application only uses single container
	environmentSlice := replicationController.ContainerSlice[0].EnvironmentSlice

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imagerecords/" + name

	imageRecordSlice := make([]ImageRecord, 0)

	_, err = restclient.RequestGetWithStructure(url, &imageRecordSlice)
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		filteredImageRecordSlice := make([]ImageRecord, 0)
		for _, imageRecord := range imageRecordSlice {
			if imageRecord.Version != oldVersion {
				for key, _ := range imageRecord.Environment {
					// clean all environment field
					imageRecord.Environment[key] = ""
				}

				for _, environment := range environmentSlice {
					// Replace environment description with the current used environment value
					imageRecord.Environment[environment.Name] = environment.Value
				}

				filteredImageRecordSlice = append(filteredImageRecordSlice, imageRecord)
			}
		}

		c.Data["name"] = name
		c.Data["imageRecordSlice"] = filteredImageRecordSlice
	}

	guimessage.OutputMessage(c.Data)
}

func (c *UpdateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/deploy/deploy/")
		return
	}
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

	_, err = restclient.RequestPutWithStructure(url, deployUpdateInput, nil)

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
