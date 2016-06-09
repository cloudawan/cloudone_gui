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

package appservice

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/dashboard"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"strings"
)

type DeployInformation struct {
	Namespace                 string
	ImageInformationName      string
	CurrentVersion            string
	CurrentVersionDescription string
	Description               string
	ReplicaAmount             int
}

type DeployClusterApplication struct {
	Name                           string
	Namespace                      string
	Size                           int
	EnvironmentSlice               []interface{}
	ServiceName                    string
	ReplicationControllerNameSlice []string
}

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)
	c.TplName = "dashboard/deploy/index.html"

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Dashboard tab menu
	user, _ := c.GetSession("user").(*rbac.User)
	c.Data["dashboardTabMenu"] = identity.GetDashboardTabMenu(user, "deploy")

	cloudoneGUIProtocol := beego.AppConfig.String("cloudoneGUIProtocol")
	cloudoneGUIHost, cloudoneGUIPort := dashboard.GetServerHostAndPortFromUserRequest(c.Ctx.Input)

	c.Data["cloudoneGUIProtocol"] = cloudoneGUIProtocol
	c.Data["cloudoneGUIHost"] = cloudoneGUIHost
	c.Data["cloudoneGUIPort"] = cloudoneGUIPort

	guimessage.OutputMessage(c.Data)
}

type DataController struct {
	beego.Controller
}

func (c *DataController) Get() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	// Json
	c.Data["json"] = make(map[string]interface{})
	c.Data["json"].(map[string]interface{})["applicationView"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["thirdpartyView"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["errorMap"] = make(map[string]interface{})

	// Application view
	applicationJsonMap := make(map[string]interface{})
	applicationJsonMap["name"] = "App View"
	applicationJsonMap["children"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["applicationView"] = append(c.Data["json"].(map[string]interface{})["applicationView"].([]interface{}), applicationJsonMap)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/"

	deployInformationSlice := make([]DeployInformation, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &deployInformationSlice, tokenHeaderMap)

	if err != nil {
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.ServeJSON()
		return
	}

	deployInformationMap := make(map[string][]DeployInformation)
	for _, deployInformation := range deployInformationSlice {
		if deployInformationMap[deployInformation.ImageInformationName] == nil {
			deployInformationMap[deployInformation.ImageInformationName] = make([]DeployInformation, 0)
		}
		deployInformationMap[deployInformation.ImageInformationName] = append(deployInformationMap[deployInformation.ImageInformationName], deployInformation)
	}

	applicationViewLeafAmount := 0
	for key, deployInformationSlice := range deployInformationMap {
		deployInformationJsonMap := make(map[string]interface{})
		deployInformationJsonMap["name"] = key
		deployInformationJsonMap["color"] = dashboard.TextColorDeployInformation
		deployInformationJsonMap["children"] = make([]interface{}, 0)

		for _, deployInformation := range deployInformationSlice {
			namespaceJsonMap := make(map[string]interface{})
			namespaceJsonMap["name"] = deployInformation.Namespace + " (" + strconv.Itoa(deployInformation.ReplicaAmount) + ")"
			namespaceJsonMap["color"] = dashboard.TextColorNamespace
			namespaceJsonMap["children"] = make([]interface{}, 0)

			versionJsonMap := make(map[string]interface{})
			versionJsonMap["name"] = deployInformation.CurrentVersion + " " + deployInformation.CurrentVersionDescription
			versionJsonMap["color"] = dashboard.TextColorVersion
			versionJsonMap["children"] = make([]interface{}, 0)

			namespaceJsonMap["children"] = append(namespaceJsonMap["children"].([]interface{}), versionJsonMap)
			deployInformationJsonMap["children"] = append(deployInformationJsonMap["children"].([]interface{}), namespaceJsonMap)
		}

		applicationJsonMap["children"] = append(applicationJsonMap["children"].([]interface{}), deployInformationJsonMap)
		applicationViewLeafAmount += 1
	}

	// Third-party view
	thirdpartyJsonMap := make(map[string]interface{})
	thirdpartyJsonMap["name"] = "3rd party View"
	thirdpartyJsonMap["children"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["thirdpartyView"] = append(c.Data["json"].(map[string]interface{})["thirdpartyView"].([]interface{}), thirdpartyJsonMap)

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deployclusterapplications/"

	deployClusterApplicationSlice := make([]DeployClusterApplication, 0)

	_, err = restclient.RequestGetWithStructure(url, &deployClusterApplicationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.ServeJSON()
		return
	}

	deployClusterApplicationMap := make(map[string][]DeployClusterApplication)
	for _, deployClusterApplication := range deployClusterApplicationSlice {
		if deployClusterApplicationMap[deployClusterApplication.Name] == nil {
			deployClusterApplicationMap[deployClusterApplication.Name] = make([]DeployClusterApplication, 0)
		}
		deployClusterApplicationMap[deployClusterApplication.Name] = append(deployClusterApplicationMap[deployClusterApplication.Name], deployClusterApplication)
	}

	thirdpartyViewLeafAmount := 0
	for deployClusterApplicationName, deployClusterApplicationSlice := range deployClusterApplicationMap {
		deployClusterApplicationJsonMap := make(map[string]interface{})
		deployClusterApplicationJsonMap["name"] = deployClusterApplicationName
		deployClusterApplicationJsonMap["color"] = dashboard.TextColorDeployClusterApplication
		deployClusterApplicationJsonMap["children"] = make([]interface{}, 0)

		for _, deployClusterApplication := range deployClusterApplicationSlice {
			namespaceJsonMap := make(map[string]interface{})
			namespaceJsonMap["name"] = deployClusterApplication.Namespace + " (" + strconv.Itoa(deployClusterApplication.Size) + ")"
			namespaceJsonMap["color"] = dashboard.TextColorNamespace
			namespaceJsonMap["children"] = make([]interface{}, 0)

			// Retrieve environment
			glusterfsEndpoint := ""
			glusterfsPathList := ""
			for _, environment := range deployClusterApplication.EnvironmentSlice {
				environmentJsonMap, _ := environment.(map[string]interface{})
				name, _ := environmentJsonMap["name"].(string)
				value, _ := environmentJsonMap["value"].(string)

				if name == "GLUSTERFS_ENDPOINTS" {
					glusterfsEndpoint = value
				} else if name == "GLUSTERFS_PATH_LIST" {
					glusterfsPathList = value
				}
			}

			// Glusterfs
			glusterfsPathSlice := make([]string, 0)
			if len(glusterfsEndpoint) > 0 && len(glusterfsPathList) > 0 {
				glusterfsPathSplits := strings.Split(glusterfsPathList, ",")
				for _, glusterfsPathSplit := range glusterfsPathSplits {
					glusterfsPathSlice = append(glusterfsPathSlice, strings.TrimSpace(glusterfsPathSplit))
				}
			}

			for _, replicationControllerName := range deployClusterApplication.ReplicationControllerNameSlice {
				replicationControllerNameJsonMap := make(map[string]interface{})
				replicationControllerNameJsonMap["name"] = replicationControllerName
				replicationControllerNameJsonMap["color"] = dashboard.TextColorReplicationController
				replicationControllerNameJsonMap["children"] = make([]interface{}, 0)

				// Glusterfs
				if len(glusterfsPathSlice) > 0 {
					nameSplits := strings.Split(replicationControllerName, "-")
					index, err := strconv.Atoi(nameSplits[len(nameSplits)-1])
					if err == nil && index < len(glusterfsPathSlice) {
						glusterfsJsonMap := make(map[string]interface{})
						glusterfsJsonMap["name"] = glusterfsEndpoint + " / " + glusterfsPathSlice[index]
						glusterfsJsonMap["color"] = dashboard.TextColorGlusterfs
						glusterfsJsonMap["children"] = make([]interface{}, 0)
						replicationControllerNameJsonMap["children"] = append(replicationControllerNameJsonMap["children"].([]interface{}), glusterfsJsonMap)
					}
				}

				namespaceJsonMap["children"] = append(namespaceJsonMap["children"].([]interface{}), replicationControllerNameJsonMap)

				thirdpartyViewLeafAmount += 1
			}

			deployClusterApplicationJsonMap["children"] = append(deployClusterApplicationJsonMap["children"].([]interface{}), namespaceJsonMap)
		}

		thirdpartyJsonMap["children"] = append(thirdpartyJsonMap["children"].([]interface{}), deployClusterApplicationJsonMap)
	}

	c.Data["json"].(map[string]interface{})["applicationViewLeafAmount"] = applicationViewLeafAmount
	c.Data["json"].(map[string]interface{})["thirdpartyViewLeafAmount"] = thirdpartyViewLeafAmount

	dashboard.RecursiveSortTheDataInGraphJsonMap(c.Data["json"].(map[string]interface{}))

	c.ServeJSON()
}
