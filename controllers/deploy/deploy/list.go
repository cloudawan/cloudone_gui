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

package deploy

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
)

type ListController struct {
	beego.Controller
}

type DeployInformation struct {
	Namespace                      string
	ImageInformationName           string
	CurrentVersion                 string
	CurrentVersionDescription      string
	Description                    string
	ReplicaAmount                  int
	AutoUpdateForNewBuild          bool
	HiddenTagGuiDeployDeployUpdate string
	HiddenTagGuiDeployDeployResize string
	HiddenTagGuiDeployDeployDelete string
}

type ByDeployInformation []DeployInformation

func (b ByDeployInformation) Len() int           { return len(b) }
func (b ByDeployInformation) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByDeployInformation) Less(i, j int) bool { return b.getIdentifier(i) < b.getIdentifier(j) }
func (b ByDeployInformation) getIdentifier(i int) string {
	return b[i].Namespace + "_" + b[i].ImageInformationName + "_" + b[i].CurrentVersion
}

func (c *ListController) Get() {
	c.TplName = "deploy/deploy/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	// Tag won't work in loop so need to be placed in data
	hiddenTagGuiDeployDeployUpdate := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/deploy/update")
	hiddenTagGuiDeployDeployResize := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/deploy/resize")
	hiddenTagGuiDeployDeployDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/deploy/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespace, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/" + namespace

	deployInformationSlice := make([]DeployInformation, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &deployInformationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		// Only show those belonging to this namespace
		filteredDeployInformationSlice := make([]DeployInformation, 0)
		for _, deployInformation := range deployInformationSlice {
			if hiddenTagGuiDeployDeployUpdate {
				deployInformation.HiddenTagGuiDeployDeployUpdate = "<div class='btn-group'>"
			} else {
				deployInformation.HiddenTagGuiDeployDeployUpdate = "<div hidden>"
			}
			if hiddenTagGuiDeployDeployResize {
				deployInformation.HiddenTagGuiDeployDeployResize = "<div class='btn-group'>"
			} else {
				deployInformation.HiddenTagGuiDeployDeployResize = "<div hidden>"
			}
			if hiddenTagGuiDeployDeployDelete {
				deployInformation.HiddenTagGuiDeployDeployDelete = "<div class='btn-group'>"
			} else {
				deployInformation.HiddenTagGuiDeployDeployDelete = "<div hidden>"
			}

			filteredDeployInformationSlice = append(filteredDeployInformationSlice, deployInformation)
		}

		sort.Sort(ByDeployInformation(filteredDeployInformationSlice))
		c.Data["deployInformationSlice"] = filteredDeployInformationSlice
	}

	guimessage.OutputMessage(c.Data)
}
