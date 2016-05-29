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

package topologytemplate

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"strings"
)

type CloneController struct {
	beego.Controller
}

func (c *CloneController) Get() {
	c.TplName = "repository/topologytemplate/clone.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	topologyName := c.GetString("name")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/topology/" + topologyName

	topology := Topology{}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &topology, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/repository/topologytemplate/list")
		return
	}

	c.Data["topology"] = topology

	if topology.LaunchSlice != nil {
		for i := 0; i < len(topology.LaunchSlice); i++ {
			if topology.LaunchSlice[i].LaunchApplication != nil {
				topology.LaunchSlice[i].Name = topology.LaunchSlice[i].LaunchApplication.ImageInformationName
			} else {
				topology.LaunchSlice[i].LaunchApplication = &LaunchApplication{}
				topology.LaunchSlice[i].HiddenTagLaunchApplication = "hidden"
				// To avoid client side input validation
				topology.LaunchSlice[i].LaunchApplication.ReplicaAmount = 1
			}
			if topology.LaunchSlice[i].LaunchClusterApplication != nil {
				topology.LaunchSlice[i].Name = topology.LaunchSlice[i].LaunchClusterApplication.Name
			} else {
				topology.LaunchSlice[i].LaunchClusterApplication = &LaunchClusterApplication{}
				topology.LaunchSlice[i].HiddenTagLaunchClusterApplication = "hidden"
				// To avoid client side input validation
				topology.LaunchSlice[i].LaunchClusterApplication.Size = 1
			}
		}
	}

	// Region
	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/nodes/topology"

	regionSlice := make([]Region, 0)

	_, err = restclient.RequestGetWithStructure(url, &regionSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/repository/topologytemplate/list")
		return
	}

	filteredRegionSlice := make([]Region, 0)
	for _, region := range regionSlice {
		if region.LocationTagged {
			filteredRegionSlice = append(filteredRegionSlice, region)
		}
	}

	for i := 0; i < len(topology.LaunchSlice); i++ {
		topology.LaunchSlice[i].RegionSlice = filteredRegionSlice
	}

	guimessage.OutputMessage(c.Data)
}

func (c *CloneController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	// Get topology
	name := c.GetString("name")
	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/topology/" + name

	topology := Topology{}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &topology, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/repository/topologytemplate/list")
		return
	}

	// Generate topology order
	inputMap := c.Input()
	launchNameSlice := make([]string, 0)
	environmentMap := make(map[string]string)
	if inputMap != nil {
		for key, _ := range inputMap {
			// Collect the clone name
			if strings.HasPrefix(key, "launchUse") {
				launchName := key[len("launchUse"):]
				if c.GetString(key) == "on" && len(launchName) > 0 {
					launchNameSlice = append(launchNameSlice, launchName)
				}
				// Collect environment
			} else if strings.HasPrefix(key, "clusterEnvironment") {
				environmentMap[key] = c.GetString(key)
			} else if strings.HasPrefix(key, "applicationEnvironment") {
				environmentMap[key] = c.GetString(key)
			}
		}
	}

	launchSlice := make([]Launch, 0)
	for _, launchName := range launchNameSlice {
		launchOrder, _ := c.GetInt("launchOrder" + launchName)

		clusterName := c.GetString("clusterName" + launchName)
		clusterSize, _ := c.GetInt("clusterSize" + launchName)

		applicationImageInformationName := c.GetString("applicationImageInformationName" + launchName)
		applicationVersion := c.GetString("applicationVersion" + launchName)
		applicationDescription := c.GetString("applicationDescription" + launchName)
		applicationReplicaAmount, _ := c.GetInt("applicationReplicaAmount" + launchName)

		clusterEnvironmentMap := make(map[string]string)
		applicationEnvironmentMap := make(map[string]string)
		for key, value := range environmentMap {
			if strings.HasPrefix(key, "clusterEnvironment"+launchName) {
				environemntKey := key[len("clusterEnvironment"+launchName):]
				clusterEnvironmentMap[environemntKey] = value
			} else if strings.HasPrefix(key, "applicationEnvironment"+launchName) {
				environemntKey := key[len("applicationEnvironment"+launchName):]
				applicationEnvironmentMap[environemntKey] = value
			}
		}

		// Location Affinity
		region := c.GetString("launchRegion" + launchName)
		zone := c.GetString("launchZone" + launchName)

		if region == "Any" {
			region = ""
		}
		if zone == "Any" {
			zone = ""
		}

		extraJsonMap := make(map[string]interface{})
		if len(region) > 0 {
			extraJsonMap["spec"] = make(map[string]interface{})
			extraJsonMap["spec"].(map[string]interface{})["template"] = make(map[string]interface{})
			extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"] = make(map[string]interface{})
			extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["nodeSelector"] = make(map[string]interface{})

			extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["nodeSelector"].(map[string]interface{})["region"] = region
			if len(zone) > 0 {
				extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["nodeSelector"].(map[string]interface{})["zone"] = zone
			}
		} else {
			extraJsonMap = nil
		}

		// Cluster or application
		if len(clusterName) > 0 {
			environmentSlice := make([]interface{}, 0)
			for key, value := range clusterEnvironmentMap {
				environmentJsonMap := make(map[string]interface{})
				environmentJsonMap["name"] = key
				environmentJsonMap["value"] = value
				environmentSlice = append(environmentSlice, environmentJsonMap)
			}

			clusterLaunch := &LaunchClusterApplication{
				clusterName,
				clusterSize,
				environmentSlice,
				extraJsonMap,
			}

			launch := Launch{
				launchOrder,
				nil,
				clusterLaunch,
				"",
				"",
				nil,
				"",
				"",
			}

			launchSlice = append(launchSlice, launch)
		} else if len(applicationImageInformationName) > 0 {
			environmentSlice := make([]ReplicationControllerContainerEnvironment, 0)
			for key, value := range applicationEnvironmentMap {
				environmentSlice = append(environmentSlice, ReplicationControllerContainerEnvironment{key, value})
			}

			oldLaunch := Launch{}
			for i := 0; i < len(topology.LaunchSlice); i++ {
				if topology.LaunchSlice[i].LaunchApplication != nil && launchName == topology.LaunchSlice[i].LaunchApplication.ImageInformationName {
					oldLaunch = topology.LaunchSlice[i]
				}
			}

			// Change the assigned Node Port to auto generated
			for i, port := range oldLaunch.LaunchApplication.PortSlice {
				if port.NodePort > 0 {
					oldLaunch.LaunchApplication.PortSlice[i].NodePort = 0
				}
			}

			launchApplication := &LaunchApplication{
				applicationImageInformationName,
				applicationVersion,
				applicationDescription,
				applicationReplicaAmount,
				oldLaunch.LaunchApplication.PortSlice,
				environmentSlice,
				oldLaunch.LaunchApplication.ResourceMap,
				extraJsonMap,
			}

			launch := Launch{
				launchOrder,
				launchApplication,
				nil,
				"",
				"",
				nil,
				"",
				"",
			}

			launchSlice = append(launchSlice, launch)
		}
	}

	sort.Sort(ByLaunch(launchSlice))

	// Clone
	namespace, _ := c.GetSession("namespace").(string)

	for _, launch := range launchSlice {
		if launch.LaunchApplication != nil {
			url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
				"/api/v1/deploys/create/" + namespace

			_, err = restclient.RequestPostWithStructure(url, launch.LaunchApplication, nil, tokenHeaderMap)

			if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
				return
			}

			if err != nil {
				// Error
				guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
				guimessage.RedirectMessage(c)
				c.Ctx.Redirect(302, "/gui/repository/topologytemplate/list")
				return
			}
		}
		if launch.LaunchClusterApplication != nil {
			url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
				"/api/v1/clusterapplications/launch/" + namespace + "/" + launch.LaunchClusterApplication.Name

			jsonMap := make(map[string]interface{})

			_, err = restclient.RequestPostWithStructure(url, launch.LaunchClusterApplication, &jsonMap, tokenHeaderMap)

			if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
				return
			}

			if err != nil {
				guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
				guimessage.RedirectMessage(c)
				c.Ctx.Redirect(302, "/gui/repository/topologytemplate/list")
				return
			}
		}
	}

	guimessage.AddSuccess("Clone from the template " + name + " to the namespace " + namespace)

	c.Ctx.Redirect(302, "/gui/repository/topologytemplate/list")

	guimessage.RedirectMessage(c)
}
