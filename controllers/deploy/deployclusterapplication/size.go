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
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
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

type SizeController struct {
	beego.Controller
}

func (c *SizeController) Get() {
	c.TplNames = "deploy/deployclusterapplication/size.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")

	namespace, _ := c.GetSession("namespace").(string)

	name := c.GetString("name")
	size := c.GetString("size")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/clusterapplications/" + name
	cluster := Cluster{}
	_, err := restclient.RequestGetWithStructure(url, &cluster)

	if err != nil {
		guimessage.AddDanger("Fail to get cluster application with error" + err.Error())
		// Redirect to list
		c.Ctx.Redirect(302, "/gui/deploy/deployclusterapplication/")

		guimessage.RedirectMessage(c)
	} else {
		c.Data["name"] = name
		c.Data["size"] = size

		// Get configured environment from the first instance
		clusterApplicationFirstInstanceName := name + "-instance-0"
		url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
			"/api/v1/replicationcontrollers/" + namespace + "/" + clusterApplicationFirstInstanceName +
			"?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort
		replicationController := ReplicationController{}
		_, err := restclient.RequestGetWithStructure(url, &replicationController)

		// Cluster application
		if err == nil {
			for _, container := range replicationController.ContainerSlice {
				for _, environment := range container.EnvironmentSlice {
					cluster.Environment[environment.Name] = environment.Value
				}
			}

			c.Data["environment"] = cluster.Environment
		}

		guimessage.OutputMessage(c.Data)
	}
}

func (c *SizeController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")

	namespace, _ := c.GetSession("namespace").(string)

	name := c.GetString("name")
	size := c.GetString("size")

	keySlice := make([]string, 0)
	inputMap := c.Input()
	if inputMap != nil {
		for key, _ := range inputMap {
			keySlice = append(keySlice, key)
		}
	}

	environmentSlice := make([]interface{}, 0)
	for _, key := range keySlice {
		value := c.GetString(key)
		if len(value) > 0 {
			environmentMap := make(map[string]string)
			environmentMap["name"] = key
			environmentMap["value"] = value
			environmentSlice = append(environmentSlice, environmentMap)
		}
	}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/deployclusterapplications/size/" + namespace + "/" + name +
		"?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort + "&size=" + size

	_, err := restclient.RequestPut(url, environmentSlice, true)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Cluster application " + name + " is resized")
	}

	c.Ctx.Redirect(302, "/gui/deploy/deployclusterapplication/")

	guimessage.RedirectMessage(c)
}