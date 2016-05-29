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

package user

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"strings"
)

type ListController struct {
	beego.Controller
}

type SimplifiedUser struct {
	Name                             string
	Disabled                         bool
	ExpiredTime                      string
	RoleNameSlice                    []string
	NamespaceSlice                   []string
	Description                      string
	HiddenTagGuiSystemRBACUserEdit   string
	HiddenTagGuiSystemRBACUserDelete string
}

type BySimplifiedUser []SimplifiedUser

func (b BySimplifiedUser) Len() int           { return len(b) }
func (b BySimplifiedUser) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b BySimplifiedUser) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "system/rbac/user/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// System RBAC tab menu
	user, _ := c.GetSession("user").(*rbac.User)
	c.Data["systemRBACTabMenu"] = identity.GetSystemRBACTabMenu(user, "user")
	// Authorization for Button
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiSystemRBACUserEdit", user, "GET", "/gui/system/rbac/user/edit")
	// Tag won't work in loop so need to be placed in data
	hasGuiSystemRBACUserEdit := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/rbac/user/edit")
	hasGuiSystemRBACUserDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/rbac/user/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/authorizations/users"

	userSlice := make([]rbac.User, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &userSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort + "/api/v1/namespaces"

		namespaceNameSlice := make([]string, 0)

		_, err = restclient.RequestGetWithStructure(url, &namespaceNameSlice, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		} else {
			simplifiedUserSlice := make([]SimplifiedUser, 0)
			for _, user := range userSlice {
				roleNameSlice := make([]string, 0)
				for _, role := range user.RoleSlice {
					roleNameSlice = append(roleNameSlice, role.Name)
				}

				namespaceSlice := make([]string, 0)
				for _, resource := range user.ResourceSlice {
					// Use component * only since they are the same in all components in the simplified version
					if resource.Component == "*" {
						if resource.Path == "*" || resource.Path == "/namespaces/" {
							namespaceSlice = make([]string, 0)
							namespaceSlice = append(namespaceSlice, "*")
							break
						} else if strings.HasPrefix(resource.Path, "/namespaces/") {
							splitSlice := strings.Split(resource.Path, "/")
							namespace := splitSlice[2]
							namespaceSlice = append(namespaceSlice, namespace)
						}
					}
				}

				expiredTimeText := ""
				if user.ExpiredTime != nil {
					expiredTimeText = user.ExpiredTime.String()
				}

				simplifiedUser := SimplifiedUser{
					user.Name,
					user.Disabled,
					expiredTimeText,
					roleNameSlice,
					namespaceSlice,
					user.Description,
					"",
					"",
				}

				if hasGuiSystemRBACUserEdit {
					simplifiedUser.HiddenTagGuiSystemRBACUserEdit = "<div class='btn-group'>"
				} else {
					simplifiedUser.HiddenTagGuiSystemRBACUserEdit = "<div hidden>"
				}
				if hasGuiSystemRBACUserDelete {
					simplifiedUser.HiddenTagGuiSystemRBACUserDelete = "<div class='btn-group'>"
				} else {
					simplifiedUser.HiddenTagGuiSystemRBACUserDelete = "<div hidden>"
				}

				simplifiedUserSlice = append(simplifiedUserSlice, simplifiedUser)
			}

			sort.Sort(BySimplifiedUser(simplifiedUserSlice))
			c.Data["simplifiedUserSlice"] = simplifiedUserSlice
		}
	}

	guimessage.OutputMessage(c.Data)
}
