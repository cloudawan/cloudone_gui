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

func (c *CreateController) Get() {
	c.TplName = "deploy/deploy/create.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imagerecords/" + name

	imageRecordSlice := make([]ImageRecord, 0)

	_, err := restclient.RequestGetWithStructure(url, &imageRecordSlice)
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["imageInformationName"] = name
		c.Data["imageRecordSlice"] = imageRecordSlice
	}

	guimessage.OutputMessage(c.Data)
}

func (c *CreateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	namespaces, _ := c.GetSession("namespace").(string)

	imageInformationName := c.GetString("imageInformationName")
	version := c.GetString("version")
	description := c.GetString("description")
	replicaAmount, _ := c.GetInt("replicaAmount")
	//portName := c.GetString("portName")
	containerPort, err := c.GetInt("containerPort")

	portName := "generated"

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

	replicationControllerContainerPortSlice := make([]ReplicationControllerContainerPort, 0)
	replicationControllerContainerPortSlice = append(replicationControllerContainerPortSlice, ReplicationControllerContainerPort{portName, containerPort})

	deployCreateInput := DeployCreateInput{
		imageInformationName,
		version,
		description,
		replicaAmount,
		replicationControllerContainerPortSlice,
		environmentSlice,
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/create/" + namespaces + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	_, err = restclient.RequestPostWithStructure(url, deployCreateInput, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Create deploy " + imageInformationName + " version " + version + " success")
	}

	c.Ctx.Redirect(302, "/gui/deploy/deploy/")

	guimessage.RedirectMessage(c)
}
