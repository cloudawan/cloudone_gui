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
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"strconv"
)

type ListController struct {
	beego.Controller
}

type DeployBlueGreen struct {
	ImageInformation                        string
	Namespace                               string
	NodePort                                int
	Description                             string
	SessionAffinity                         string
	NodePortDisplay                         string
	HiddenTagGuiDeployDeployBlueGreenSelect string
	HiddenTagGuiDeployDeployBlueGreenDelete string
}

type ByDeployBlueGreen []DeployBlueGreen

func (b ByDeployBlueGreen) Len() int           { return len(b) }
func (b ByDeployBlueGreen) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByDeployBlueGreen) Less(i, j int) bool { return b.getIdentifier(i) < b.getIdentifier(j) }
func (b ByDeployBlueGreen) getIdentifier(i int) string {
	return b[i].Namespace + "_" + b[i].ImageInformation
}

func (c *ListController) Get() {
	c.TplName = "deploy/deploybluegreen/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	// Tag won't work in loop so need to be placed in data
	hasGuiDeployDeployBlueGreenSelect := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/deploybluegreen/select")
	hasGuiDeployDeployBlueGreenDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/deploybluegreen/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploybluegreens/"

	deployBlueGreenSlice := make([]DeployBlueGreen, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &deployBlueGreenSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(deployBlueGreenSlice); i++ {
			if deployBlueGreenSlice[i].NodePort == 0 {
				deployBlueGreenSlice[i].NodePortDisplay = "Auto-generated"
			} else {
				deployBlueGreenSlice[i].NodePortDisplay = strconv.Itoa(deployBlueGreenSlice[i].NodePort)
			}

			if hasGuiDeployDeployBlueGreenSelect {
				deployBlueGreenSlice[i].HiddenTagGuiDeployDeployBlueGreenSelect = "<div class='btn-group'>"
			} else {
				deployBlueGreenSlice[i].HiddenTagGuiDeployDeployBlueGreenSelect = "<div hidden>"
			}
			if hasGuiDeployDeployBlueGreenDelete {
				deployBlueGreenSlice[i].HiddenTagGuiDeployDeployBlueGreenDelete = "<div class='btn-group'>"
			} else {
				deployBlueGreenSlice[i].HiddenTagGuiDeployDeployBlueGreenDelete = "<div hidden>"
			}
		}

		sort.Sort(ByDeployBlueGreen(deployBlueGreenSlice))
		c.Data["deployBlueGreenSlice"] = deployBlueGreenSlice
	}

	guimessage.OutputMessage(c.Data)
}
