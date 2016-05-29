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

package server

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

type PrivateRegistry struct {
	Name                                            string
	Host                                            string
	Port                                            int
	HiddenTagGuiSystemPrivateRegistryRepositoryList string
	HiddenTagGuiSystemPrivateRegistryServerEdit     string
	HiddenTagGuiSystemPrivateRegistryServerDelete   string
}

type ByPrivateRegistry []PrivateRegistry

func (b ByPrivateRegistry) Len() int           { return len(b) }
func (b ByPrivateRegistry) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPrivateRegistry) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "system/privateregistry/server/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiSystemPrivateRegistryServerEdit", user, "GET", "/gui/system/privateregistry/server/edit")
	// Tag won't work in loop so need to be placed in data
	hasGuiSystemPrivateRegistryRepositoryList := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/privateregistry/repository/list")
	hasGuiSystemPrivateRegistryServerEdit := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/privateregistry/server/edit")
	hasGuiSystemPrivateRegistryServerDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/privateregistry/server/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/privateregistries/servers/"

	privateRegistrySlice := make([]PrivateRegistry, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &privateRegistrySlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		for i := 0; i < len(privateRegistrySlice); i++ {
			if hasGuiSystemPrivateRegistryRepositoryList {
				privateRegistrySlice[i].HiddenTagGuiSystemPrivateRegistryRepositoryList = "<div class='btn-group'>"
			} else {
				privateRegistrySlice[i].HiddenTagGuiSystemPrivateRegistryRepositoryList = "<div hidden>"
			}
			if hasGuiSystemPrivateRegistryServerEdit {
				privateRegistrySlice[i].HiddenTagGuiSystemPrivateRegistryServerEdit = "<div class='btn-group'>"
			} else {
				privateRegistrySlice[i].HiddenTagGuiSystemPrivateRegistryServerEdit = "<div hidden>"
			}
			if hasGuiSystemPrivateRegistryServerDelete {
				privateRegistrySlice[i].HiddenTagGuiSystemPrivateRegistryServerDelete = "<div class='btn-group'>"
			} else {
				privateRegistrySlice[i].HiddenTagGuiSystemPrivateRegistryServerDelete = "<div hidden>"
			}
		}

		sort.Sort(ByPrivateRegistry(privateRegistrySlice))
		c.Data["privateRegistrySlice"] = privateRegistrySlice
	}

	guimessage.OutputMessage(c.Data)
}
