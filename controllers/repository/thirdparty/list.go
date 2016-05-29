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

package thirdparty

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

type ThirdPartyApplication struct {
	Name                                   string
	Description                            string
	HiddenTagGuiRepositoryThirdPartyLaunch string
	HiddenTagGuiRepositoryThirdPartyEdit   string
	HiddenTagGuiRepositoryThirdPartyDelete string
}

type ByThirdPartyApplication []ThirdPartyApplication

func (b ByThirdPartyApplication) Len() int           { return len(b) }
func (b ByThirdPartyApplication) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByThirdPartyApplication) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "repository/thirdparty/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiRepositoryThirdPartyEdit", user, "GET", "/gui/repository/thirdparty/edit")
	// Tag won't work in loop so need to be placed in data
	hasGuiRepositoryThirdPartyLaunch := user.HasPermission(identity.GetConponentName(), "GET", "/gui/repository/thirdparty/launch")
	hasGuiRepositoryThirdPartyEdit := user.HasPermission(identity.GetConponentName(), "GET", "/gui/repository/thirdparty/edit")
	hasGuiRepositoryThirdPartyDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/repository/thirdparty/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/clusterapplications/"

	thirdPartyApplicationSlice := make([]ThirdPartyApplication, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &thirdPartyApplicationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		for i := 0; i < len(thirdPartyApplicationSlice); i++ {
			if hasGuiRepositoryThirdPartyLaunch {
				thirdPartyApplicationSlice[i].HiddenTagGuiRepositoryThirdPartyLaunch = "<div class='btn-group'>"
			} else {
				thirdPartyApplicationSlice[i].HiddenTagGuiRepositoryThirdPartyLaunch = "<div hidden>"
			}
			if hasGuiRepositoryThirdPartyEdit {
				thirdPartyApplicationSlice[i].HiddenTagGuiRepositoryThirdPartyEdit = "<div class='btn-group'>"
			} else {
				thirdPartyApplicationSlice[i].HiddenTagGuiRepositoryThirdPartyEdit = "<div hidden>"
			}
			if hasGuiRepositoryThirdPartyDelete {
				thirdPartyApplicationSlice[i].HiddenTagGuiRepositoryThirdPartyDelete = "<div class='btn-group'>"
			} else {
				thirdPartyApplicationSlice[i].HiddenTagGuiRepositoryThirdPartyDelete = "<div hidden>"
			}
		}

		sort.Sort(ByThirdPartyApplication(thirdPartyApplicationSlice))
		c.Data["thirdPartyApplicationSlice"] = thirdPartyApplicationSlice
	}

	guimessage.OutputMessage(c.Data)
}
