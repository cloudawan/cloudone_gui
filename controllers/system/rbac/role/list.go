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

package role

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

type SimplifiedRole struct {
	Name                             string
	Description                      string
	HiddenTagGuiSystemRBACRoleEdit   string
	HiddenTagGuiSystemRBACRoleDelete string
}

type BySimplifiedRole []SimplifiedRole

func (b BySimplifiedRole) Len() int           { return len(b) }
func (b BySimplifiedRole) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b BySimplifiedRole) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "system/rbac/role/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// System RBAC tab menu
	user, _ := c.GetSession("user").(*rbac.User)
	c.Data["systemRBACTabMenu"] = identity.GetSystemRBACTabMenu(user, "role")
	// Authorization for Button
	identity.SetPriviledgeHiddenTag(c.Data, "hiddenTagGuiSystemRBACRoleEdit", user, "GET", "/gui/system/rbac/role/edit")
	// Tag won't work in loop so need to be placed in data
	hasGuiSystemRBACRoleEdit := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/rbac/role/edit")
	hasGuiSystemRBACRoleDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/rbac/role/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/authorizations/roles"

	simplifiedRoleSlice := make([]SimplifiedRole, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &simplifiedRoleSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(simplifiedRoleSlice); i++ {
			if hasGuiSystemRBACRoleEdit {
				simplifiedRoleSlice[i].HiddenTagGuiSystemRBACRoleEdit = "<div class='btn-group'>"
			} else {
				simplifiedRoleSlice[i].HiddenTagGuiSystemRBACRoleEdit = "<div hidden>"
			}
			if hasGuiSystemRBACRoleDelete {
				simplifiedRoleSlice[i].HiddenTagGuiSystemRBACRoleDelete = "<div class='btn-group'>"
			} else {
				simplifiedRoleSlice[i].HiddenTagGuiSystemRBACRoleDelete = "<div hidden>"
			}
		}

		sort.Sort(BySimplifiedRole(simplifiedRoleSlice))
		c.Data["simplifiedRoleSlice"] = simplifiedRoleSlice
	}

	guimessage.OutputMessage(c.Data)
}
