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

package topology

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type ReplicationControllerAndRelatedPod struct {
	Name               string
	Namespace          string
	ReplicaAmount      int
	AliveReplicaAmount int
	Selector           map[string]string
	Label              map[string]string
	PodSlice           []Pod
}

type Pod struct {
	Name           string
	Namespace      string
	HostIP         string
	PodIP          string
	ContainerSlice []PodContainer
}

type PodContainer struct {
	Name      string
	Image     string
	PortSlice []PodContainerPort
}

type PodContainerPort struct {
	Name          string
	ContainerPort int
	Protocol      string
}

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	c.TplName = "dashboard/topology/index.html"

	cloudoneGUIProtocol := beego.AppConfig.String("cloudoneGUIProtocol")
	cloudoneGUIHost := c.Ctx.Input.Host()
	cloudoneGUIPort := c.Ctx.Input.Port()

	c.Data["cloudoneGUIProtocol"] = cloudoneGUIProtocol
	c.Data["cloudoneGUIHost"] = cloudoneGUIHost
	c.Data["cloudoneGUIPort"] = cloudoneGUIPort

	guimessage.OutputMessage(c.Data)
}

const (
	allKeyword = "All"
)

type DataController struct {
	beego.Controller
}

func (c *DataController) Get() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	scope := c.GetString("scope")

	namespaceSlice := make([]string, 0)
	if scope == allKeyword {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/namespaces/" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

		_, err := restclient.RequestGetWithStructure(url, &namespaceSlice)
		if err != nil {
			c.Data["json"].(map[string]interface{})["error"] = err.Error()
			c.ServeJSON()
			return
		}
	} else {
		namespace, _ := c.GetSession("namespace").(string)
		namespaceSlice = append(namespaceSlice, namespace)
	}

	// Json
	c.Data["json"] = make(map[string]interface{})
	c.Data["json"].(map[string]interface{})["logicalView"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["physicalView"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["errorMap"] = make(map[string]interface{})

	// Logical view
	logicalTopologyJsonMap := make(map[string]interface{})
	logicalTopologyJsonMap["name"] = "Logical View"
	logicalTopologyJsonMap["children"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["logicalView"] = append(c.Data["json"].(map[string]interface{})["logicalView"].([]interface{}), logicalTopologyJsonMap)

	// Physical view
	physicalTopologyJsonMap := make(map[string]interface{})
	physicalTopologyJsonMap["name"] = "Physical View"
	physicalTopologyJsonMap["children"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["physicalView"] = append(c.Data["json"].(map[string]interface{})["physicalView"].([]interface{}), physicalTopologyJsonMap)

	//namespaceSlice := make([]string, 0)
	//namespaceSlice = append(namespaceSlice, namespace)
	for _, namespace := range namespaceSlice {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/replicationcontrollers/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

		replicationControllerAndRelatedPodSlice := make([]ReplicationControllerAndRelatedPod, 0)
		_, err := restclient.RequestGetWithStructure(url, &replicationControllerAndRelatedPodSlice)
		if err != nil {
			c.Data["json"].(map[string]interface{})["error"] = "Get replication controller data error"
			c.Data["json"].(map[string]interface{})["errorMap"].(map[string]interface{})[namespace] = err.Error()
		} else {
			// Logical view
			logicalTopologyNamespaceJsonMap := make(map[string]interface{})
			logicalTopologyNamespaceJsonMap["name"] = namespace
			logicalTopologyNamespaceJsonMap["children"] = make([]interface{}, 0)
			for _, replicationControllerAndRelatedPod := range replicationControllerAndRelatedPodSlice {
				replicationControllerJsonMap := make(map[string]interface{})
				replicationControllerJsonMap["name"] = replicationControllerAndRelatedPod.Name
				replicationControllerJsonMap["children"] = make([]interface{}, 0)
				for _, pod := range replicationControllerAndRelatedPod.PodSlice {
					podJsonMap := make(map[string]interface{})
					podJsonMap["name"] = pod.Name
					podJsonMap["children"] = make([]interface{}, 0)
					for _, container := range pod.ContainerSlice {
						containerJsonMap := make(map[string]interface{})
						containerJsonMap["name"] = container.Name
						containerJsonMap["children"] = make([]interface{}, 0)
						for _, port := range container.PortSlice {
							portJsonMap := make(map[string]interface{})
							portJsonMap["name"] = port.Name + " " + port.Protocol + " " + strconv.Itoa(port.ContainerPort)
							containerJsonMap["children"] = append(containerJsonMap["children"].([]interface{}), portJsonMap)
						}
						podJsonMap["children"] = append(podJsonMap["children"].([]interface{}), containerJsonMap)
					}
					replicationControllerJsonMap["children"] = append(replicationControllerJsonMap["children"].([]interface{}), podJsonMap)
				}
				logicalTopologyNamespaceJsonMap["children"] = append(logicalTopologyNamespaceJsonMap["children"].([]interface{}), replicationControllerJsonMap)
			}
			// Collect all nodes
			nodeMap := make(map[string]bool)
			for _, replicationControllerAndRelatedPod := range replicationControllerAndRelatedPodSlice {
				for _, pod := range replicationControllerAndRelatedPod.PodSlice {
					nodeMap[pod.HostIP] = true
				}
			}
			// Physical view
			for node, _ := range nodeMap {
				exist := false
				nodeJsonMap := make(map[string]interface{})
				for _, existingNodeJsonMap := range physicalTopologyJsonMap["children"].([]interface{}) {
					if existingNodeJsonMap.(map[string]interface{})["name"] == node {
						nodeJsonMap = existingNodeJsonMap.(map[string]interface{})
						exist = true
						break
					}
				}

				if exist == false {
					nodeJsonMap["name"] = node
					nodeJsonMap["children"] = make([]interface{}, 0)
				}

				for _, replicationControllerAndRelatedPod := range replicationControllerAndRelatedPodSlice {
					for _, pod := range replicationControllerAndRelatedPod.PodSlice {
						if pod.HostIP == node {
							podJsonMap := make(map[string]interface{})
							podJsonMap["name"] = pod.Name
							podJsonMap["children"] = make([]interface{}, 0)
							for _, container := range pod.ContainerSlice {
								containerJsonMap := make(map[string]interface{})
								containerJsonMap["name"] = container.Name
								containerJsonMap["children"] = make([]interface{}, 0)
								for _, port := range container.PortSlice {
									portJsonMap := make(map[string]interface{})
									portJsonMap["name"] = port.Name + " " + port.Protocol + " " + strconv.Itoa(port.ContainerPort)
									containerJsonMap["children"] = append(containerJsonMap["children"].([]interface{}), portJsonMap)
								}
								podJsonMap["children"] = append(podJsonMap["children"].([]interface{}), containerJsonMap)
							}
							nodeJsonMap["children"] = append(nodeJsonMap["children"].([]interface{}), podJsonMap)
						}
					}
				}
				if exist == false {
					physicalTopologyJsonMap["children"] = append(physicalTopologyJsonMap["children"].([]interface{}), nodeJsonMap)
				}
			}

			logicalTopologyJsonMap["children"] = append(logicalTopologyJsonMap["children"].([]interface{}), logicalTopologyNamespaceJsonMap)
		}
	}

	c.ServeJSON()
}
