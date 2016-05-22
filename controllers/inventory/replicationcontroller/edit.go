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
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplName = "inventory/replicationcontroller/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	replicationcontroller := c.GetString("replicationcontroller")
	if replicationcontroller == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Replication Controller"
		c.Data["replicationControllerName"] = ""
	} else {
		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Replication Controller"
		c.Data["replicationControllerName"] = replicationcontroller
	}

	guimessage.OutputMessage(c.Data)
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

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespace, _ := c.GetSession("namespace").(string)

	selectorName := c.GetString("name")
	replicaAmount, _ := c.GetInt("replicaAmount")
	image := c.GetString("image")
	containerPort, err := c.GetInt("containerPort")

	version := ""
	name := selectorName + version

	portName := "generated"

	replicationControllerContainerPortSlice := make([]ReplicationControllerContainerPort, 0)
	replicationControllerContainerPortSlice = append(replicationControllerContainerPortSlice, ReplicationControllerContainerPort{portName, containerPort})
	replicationControllerContainerSlice := make([]ReplicationControllerContainer, 0)
	replicationControllerContainerSlice = append(replicationControllerContainerSlice, ReplicationControllerContainer{name, image, replicationControllerContainerPortSlice})
	replicationController := ReplicationController{
		name,
		replicaAmount,
		ReplicationControllerSelector{
			selectorName,
			version,
		},
		ReplicationControllerLabel{
			name,
		},
		replicationControllerContainerSlice}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/replicationcontrollers/" + namespace

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPostWithStructure(url, replicationController, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Replication Controller " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/inventory/replicationcontroller/list")

	guimessage.RedirectMessage(c)
}
