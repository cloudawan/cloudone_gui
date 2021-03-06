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
	"github.com/cloudawan/cloudone_gui/controllers/identity"
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

// @Title Get
// @Description get topology
// @Success 200 (string) {}
// @router / [get]
func (c *IndexController) Get() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespace, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/replicationcontrollers/" + namespace

	replicationControllerAndRelatedPodSlice := make([]ReplicationControllerAndRelatedPod, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &replicationControllerAndRelatedPodSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		c.Data["json"] = `{"error": "` + err.Error() + `"}`
	} else {
		// Logical view
		logicalTopologyJsonMap := make(map[string]interface{})
		logicalTopologyJsonMap["name"] = "Logical View"
		logicalTopologyJsonMap["children"] = make([]interface{}, 0)
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
			logicalTopologyJsonMap["children"] = append(logicalTopologyJsonMap["children"].([]interface{}), replicationControllerJsonMap)
		}
		// Collect all nodes
		nodeMap := make(map[string]bool)
		for _, replicationControllerAndRelatedPod := range replicationControllerAndRelatedPodSlice {
			for _, pod := range replicationControllerAndRelatedPod.PodSlice {
				nodeMap[pod.HostIP] = true
			}
		}
		// Physical view
		physicalTopologyJsonMap := make(map[string]interface{})
		physicalTopologyJsonMap["name"] = "Physical View"
		physicalTopologyJsonMap["children"] = make([]interface{}, 0)
		for node, _ := range nodeMap {
			nodeJsonMap := make(map[string]interface{})
			nodeJsonMap["name"] = node
			nodeJsonMap["children"] = make([]interface{}, 0)
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
			physicalTopologyJsonMap["children"] = append(physicalTopologyJsonMap["children"].([]interface{}), nodeJsonMap)
		}
		// Json
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["logicalView"] = make([]interface{}, 0)
		c.Data["json"].(map[string]interface{})["logicalView"] = append(c.Data["json"].(map[string]interface{})["logicalView"].([]interface{}), logicalTopologyJsonMap)
		c.Data["json"].(map[string]interface{})["physicalView"] = make([]interface{}, 0)
		c.Data["json"].(map[string]interface{})["physicalView"] = append(c.Data["json"].(map[string]interface{})["physicalView"].([]interface{}), physicalTopologyJsonMap)
	}

	c.ServeJSON()
}
