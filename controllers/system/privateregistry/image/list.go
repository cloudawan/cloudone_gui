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

package image

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

type Image struct {
	Tag                                          string
	Server                                       string
	Repository                                   string
	HiddenTagGuiSystemPrivateRegistryImageDelete string
}

type ByImage []Image

func (b ByImage) Len() int           { return len(b) }
func (b ByImage) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByImage) Less(i, j int) bool { return b[i].Tag < b[j].Tag }

func (c *ListController) Get() {
	c.TplName = "system/privateregistry/image/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	serverName := c.GetString("serverName")
	repositoryName := c.GetString("repositoryName")

	c.Data["serverName"] = serverName

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiSystemPrivateRegistryRepositoryList", user, "GET", "/gui/system/privateregistry/repository/list")
	// Tag won't work in loop so need to be placed in data
	hasGuiSystemPrivateRegistryImageDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/privateregistry/image/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/privateregistries/servers/" + serverName + "/repositories/" + repositoryName + "/tags/"

	tagSlice := make([]string, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &tagSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		imageSlice := make([]Image, 0)
		for _, tag := range tagSlice {
			imageSlice = append(imageSlice, Image{tag, serverName, repositoryName, ""})
		}

		for i := 0; i < len(imageSlice); i++ {
			if hasGuiSystemPrivateRegistryImageDelete {
				imageSlice[i].HiddenTagGuiSystemPrivateRegistryImageDelete = "<div class='btn-group'>"
			} else {
				imageSlice[i].HiddenTagGuiSystemPrivateRegistryImageDelete = "<div hidden>"
			}
		}

		sort.Sort(ByImage(imageSlice))
		c.Data["imageSlice"] = imageSlice
	}

	guimessage.OutputMessage(c.Data)
}
