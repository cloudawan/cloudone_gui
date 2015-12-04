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

package deploybluegreen

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type SelectController struct {
	beego.Controller
}

func (c *SelectController) Get() {
	c.TplNames = "deploy/deploybluegreen/select.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	imageInformation := c.GetString("imageInformation")
	currentEnvironment := c.GetString("currentEnvironment")
	nodePort := c.GetString("nodePort")
	description := c.GetString("description")

	c.Data["actionButtonValue"] = "Update"
	c.Data["pageHeader"] = "Update Blue Green Deployment"
	c.Data["imageInformation"] = imageInformation

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploybluegreens/deployable/" + imageInformation + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	namespaceSlice := make([]string, 0)
	_, err := restclient.RequestGetWithStructure(url, &namespaceSlice)
	if err != nil {
		guimessage.AddDanger("Fail to get deployable namespace")
	} else {
		c.Data["namespaceSlice"] = namespaceSlice
		c.Data["currentEnvironment"] = currentEnvironment
		c.Data["nodePort"] = nodePort
		c.Data["description"] = description
	}

	guimessage.OutputMessage(c.Data)
}

func (c *SelectController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, _ := configuration.GetAvailableKubeapiHostAndPort()

	imageInformation := c.GetString("imageInformation")
	namespace := c.GetString("namespace")
	nodePort, _ := c.GetInt("nodePort")
	description := c.GetString("description")
	sessionAffinity := c.GetString("sessionAffinity")

	deployBlueGreen := DeployBlueGreen{
		imageInformation,
		namespace,
		nodePort,
		description,
		sessionAffinity,
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploybluegreens/" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	_, err := restclient.RequestPutWithStructure(url, deployBlueGreen, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Create blue green deployment " + imageInformation + " success")
	}

	c.Ctx.Redirect(302, "/gui/deploy/deploybluegreen/")

	guimessage.RedirectMessage(c)
}
