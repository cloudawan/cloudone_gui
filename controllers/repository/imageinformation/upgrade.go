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

package imageinformation

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type UpgradeController struct {
	beego.Controller
}

type DeployUpgradeInput struct {
	ImageInformationName string
	Description          string
}

func (c *UpgradeController) Get() {
	c.TplName = "repository/imageinformation/upgrade.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	name := c.GetString("name")
	c.Data["name"] = name

	guimessage.OutputMessage(c.Data)
}

func (c *UpgradeController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/repository/imageinformation/list")
		return
	}

	imageInformationName := c.GetString("name")
	description := c.GetString("description")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imageinformations/upgrade/?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	deployUpgradeInput := DeployUpgradeInput{imageInformationName, description}

	resultJsonMap := make(map[string]interface{})

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPutWithStructure(url, deployUpgradeInput, &resultJsonMap, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess(imageInformationName + " is launched")
	}

	c.Ctx.Redirect(302, "/gui/repository/imageinformation/log?imageInformation="+imageInformationName)

	guimessage.RedirectMessage(c)
}
