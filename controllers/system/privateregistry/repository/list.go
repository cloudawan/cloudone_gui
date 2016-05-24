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

package repository

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

type Repository struct {
	Name                                              string
	Server                                            string
	HiddenTagGuiSystemPrivateRegistryRepositoryList   string
	HiddenTagGuiSystemPrivateRegistryRepositoryDelete string
}

type ByRepository []Repository

func (b ByRepository) Len() int           { return len(b) }
func (b ByRepository) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByRepository) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "system/privateregistry/repository/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	serverName := c.GetString("serverName")

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiSystemPrivateRegistryServerList", user, "GET", "/gui/system/privateregistry/server/list")
	// Tag won't work in loop so need to be placed in data
	hasGuiSystemPrivateRegistryImageList := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/privateregistry/image/list")
	hasGuiSystemPrivateRegistryRepositoryDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/privateregistry/repository/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/privateregistries/servers/" + serverName + "/repositories/"

	nameSlice := make([]string, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &nameSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {

		repositorySlice := make([]Repository, 0)
		for _, name := range nameSlice {
			repositorySlice = append(repositorySlice, Repository{name, serverName, "", ""})
		}

		for i := 0; i < len(repositorySlice); i++ {
			if hasGuiSystemPrivateRegistryImageList {
				repositorySlice[i].HiddenTagGuiSystemPrivateRegistryRepositoryList = "<div class='btn-group'>"
			} else {
				repositorySlice[i].HiddenTagGuiSystemPrivateRegistryRepositoryList = "<div hidden>"
			}
			if hasGuiSystemPrivateRegistryRepositoryDelete {
				repositorySlice[i].HiddenTagGuiSystemPrivateRegistryRepositoryDelete = "<div class='btn-group'>"
			} else {
				repositorySlice[i].HiddenTagGuiSystemPrivateRegistryRepositoryDelete = "<div hidden>"
			}
		}

		sort.Sort(ByRepository(repositorySlice))
		c.Data["repositorySlice"] = repositorySlice
	}

	guimessage.OutputMessage(c.Data)
}
