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

package bluegreen

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/dashboard"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type DeployBlueGreen struct {
	ImageInformation string
	Namespace        string
	NodePort         int
	Description      string
	SessionAffinity  string
}

type DeployInformation struct {
	Namespace                 string
	ImageInformationName      string
	CurrentVersion            string
	CurrentVersionDescription string
	Description               string
}

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)
	c.TplName = "dashboard/bluegreen/index.html"

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Dashboard tab menu
	user, _ := c.GetSession("user").(*rbac.User)
	c.Data["dashboardTabMenu"] = identity.GetDashboardTabMenu(user, "bluegreen")

	cloudoneGUIProtocol := beego.AppConfig.String("cloudoneGUIProtocol")
	cloudoneGUIHost := c.Ctx.Input.Host()
	cloudoneGUIPort := c.Ctx.Input.Port()

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
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.ServeJSON()
		return
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploybluegreens/"

	deployBlueGreenSlice := make([]DeployBlueGreen, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestGetWithStructure(url, &deployBlueGreenSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.ServeJSON()
		return
	}

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/"

	deployInformationSlice := make([]DeployInformation, 0)

	_, err = restclient.RequestGetWithStructure(url, &deployInformationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.ServeJSON()
		return
	}

	deployInformationMap := make(map[string]DeployInformation)
	for _, deployInformation := range deployInformationSlice {
		deployInformationMap[deployInformation.ImageInformationName+"-"+deployInformation.Namespace] = deployInformation
	}

	// Json
	c.Data["json"] = make(map[string]interface{})
	c.Data["json"].(map[string]interface{})["bluegreenView"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["errorMap"] = make(map[string]interface{})

	bluegreenJsonMap := make(map[string]interface{})
	bluegreenJsonMap["name"] = "B/G View"
	bluegreenJsonMap["children"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["bluegreenView"] = append(c.Data["json"].(map[string]interface{})["bluegreenView"].([]interface{}), bluegreenJsonMap)

	leafAmount := 0
	for _, deployBlueGreen := range deployBlueGreenSlice {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/deploybluegreens/deployable/" + deployBlueGreen.ImageInformation + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

		namespaceSlice := make([]string, 0)

		_, err := restclient.RequestGetWithStructure(url, &namespaceSlice, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			c.Data["json"].(map[string]interface{})["error"] = "Get deployable namespace data error"
			c.Data["json"].(map[string]interface{})["errorMap"].(map[string]interface{})[deployBlueGreen.ImageInformation] = err.Error()
		} else {
			blueGreenJsonMap := make(map[string]interface{})
			blueGreenJsonMap["name"] = deployBlueGreen.ImageInformation
			blueGreenJsonMap["color"] = dashboard.TextColorBlueGreen
			blueGreenJsonMap["children"] = make([]interface{}, 0)
			for _, namespace := range namespaceSlice {
				deployInformation := deployInformationMap[deployBlueGreen.ImageInformation+"-"+namespace]
				if deployInformation.CurrentVersion == "" {
					c.Data["json"].(map[string]interface{})["error"] = "Get deployInformation data error"
					errorMessage, _ := c.Data["json"].(map[string]interface{})["errorMap"].(map[string]interface{})[deployBlueGreen.ImageInformation].(string)
					c.Data["json"].(map[string]interface{})["errorMap"].(map[string]interface{})[deployBlueGreen.ImageInformation] = errorMessage + " Get deployInformation data from namespace " + namespace + " error"
				} else {
					namespaceJsonMap := make(map[string]interface{})
					namespaceJsonMap["name"] = namespace
					namespaceJsonMap["color"] = dashboard.TextColorNamespace
					namespaceJsonMap["children"] = make([]interface{}, 0)
					deployInformationJsonMap := make(map[string]interface{})
					deployInformationJsonMap["name"] = deployInformation.CurrentVersion + " " + deployInformation.CurrentVersionDescription
					deployInformationJsonMap["color"] = dashboard.TextColorDeployInformation
					if namespace == deployBlueGreen.Namespace {
						deployInformationJsonMap["children"] = make([]interface{}, 0)
						nodePortJsonMap := make(map[string]interface{})
						nodePortJsonMap["name"] = "bg_" + deployBlueGreen.ImageInformation + " " + strconv.Itoa(deployBlueGreen.NodePort)
						nodePortJsonMap["color"] = dashboard.TextColorPort
						deployInformationJsonMap["children"] = append(deployInformationJsonMap["children"].([]interface{}), nodePortJsonMap)
					}
					namespaceJsonMap["children"] = append(namespaceJsonMap["children"].([]interface{}), deployInformationJsonMap)
					blueGreenJsonMap["children"] = append(blueGreenJsonMap["children"].([]interface{}), namespaceJsonMap)
				}
				leafAmount++
			}
			bluegreenJsonMap["children"] = append(bluegreenJsonMap["children"].([]interface{}), blueGreenJsonMap)
		}
	}
	c.Data["json"].(map[string]interface{})["leafAmount"] = leafAmount

	dashboard.RecursiveSortTheDataInGraphJsonMap(c.Data["json"].(map[string]interface{}))

	c.ServeJSON()
}
