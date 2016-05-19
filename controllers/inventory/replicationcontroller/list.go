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
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type ListController struct {
	beego.Controller
}

type ReplicationControllerAndRelatedPod struct {
	Name                                                     string
	Namespace                                                string
	ReplicaAmount                                            int
	AliveReplicaAmount                                       int
	Selector                                                 map[string]string
	Label                                                    map[string]string
	PodSlice                                                 []Pod
	Display                                                  string
	HiddenTagGuiInventoryReplicationControllerSize           string
	HiddenTagGuiInventoryReplicationControllerDelete         string
	HiddenTagGuiInventoryReplicationControllerPodlog         string
	HiddenTagGuiInventoryReplicationControllerPodDelete      string
	HiddenTagGuiInventoryReplicationControllerDockerterminal string
}

type Pod struct {
	Name           string
	Namespace      string
	HostIP         string
	PodIP          string
	Phase          string
	Age            string
	ContainerSlice []PodContainer
}

type PodContainer struct {
	Name         string
	Image        string
	ContainerID  string
	RestartCount int
	Ready        bool
	PortSlice    []PodContainerPort
}

type PodContainerPort struct {
	Name          string
	ContainerPort int
	Protocol      string
}

var displayMap map[string]string = map[string]string{
	"kube-dns-v6":           "disabled",
	"private-registry":      "disabled",
	"kubernetes-management": "disabled",
}

func (c *ListController) Get() {
	c.TplName = "inventory/replicationcontroller/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPriviledgeHiddenTag(c.Data, "hiddenTagGuiInventoryReplicationControllerEdit", user, "GET", "/gui/inventory/replicationcontroller/edit")
	// Tag won't work in loop so need to be placed in data
	hasGuiInventoryReplicationControllerSize := user.HasPermission(identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/size")
	hasGuiInventoryReplicationControllerDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/delete")
	hasGuiInventoryReplicationControllerPodlog := user.HasPermission(identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/pod/log")
	hasGuiInventoryReplicationControllerPodDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/pod/delete")
	hasGuiInventoryReplicationControllerDockerterminal := user.HasPermission(identity.GetConponentName(), "GET", "/gui/inventory/replicationcontroller/dockerterminal")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.OutputMessage(c.Data)
		return
	}

	namespace, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/replicationcontrollers/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	replicationControllerAndRelatedPodSlice := make([]ReplicationControllerAndRelatedPod, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestGetWithStructure(url, &replicationControllerAndRelatedPodSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(replicationControllerAndRelatedPodSlice); i++ {
			replicationControllerAndRelatedPodSlice[i].Display =
				displayMap[replicationControllerAndRelatedPodSlice[i].Name]

			if hasGuiInventoryReplicationControllerSize {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerSize = "<div class='btn-group'>"
			} else {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerSize = "<div hidden>"
			}
			if hasGuiInventoryReplicationControllerDelete {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerDelete = "<div class='btn-group'>"
			} else {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerDelete = "<div hidden>"
			}
			if hasGuiInventoryReplicationControllerPodlog {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerPodlog = "<div class='btn-group'>"
			} else {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerPodlog = "<div hidden>"
			}
			if hasGuiInventoryReplicationControllerPodDelete {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerPodDelete = "<div class='btn-group'>"
			} else {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerPodDelete = "<div hidden>"
			}
			if hasGuiInventoryReplicationControllerDockerterminal {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerDockerterminal = "<div class='btn-group'>"
			} else {
				replicationControllerAndRelatedPodSlice[i].HiddenTagGuiInventoryReplicationControllerDockerterminal = "<div hidden>"
			}
		}
		c.Data["replicationControllerAndRelatedPodSlice"] = replicationControllerAndRelatedPodSlice
	}

	guimessage.OutputMessage(c.Data)
}
