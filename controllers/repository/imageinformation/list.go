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
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
)

type ListController struct {
	beego.Controller
}

type ImageInformation struct {
	Name                                          string
	Kind                                          string
	Description                                   string
	CurrentVersion                                string
	BuildParameter                                map[string]string
	HiddenTagGuiRepositoryImageRecordList         string
	HiddenTagGuiRepositoryImageInformationUpgrade string
	HiddenTagGuiRepositoryImageInformationLog     string
	HiddenTagGuiDeployDeployCreate                string
	HiddenTagGuiDeployDeployBlueGreenSelect       string
	HiddenTagGuiRepositoryImageInformationDelete  string
}

type ByImageInformation []ImageInformation

func (b ByImageInformation) Len() int           { return len(b) }
func (b ByImageInformation) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByImageInformation) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "repository/imageinformation/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPriviledgeHiddenTag(c.Data, "hiddenTagGuiRepositoryImageInformationCreate", user, "GET", "/gui/repository/imageinformation/create")
	// Tag won't work in loop so need to be placed in data
	hasGuiRepositoryImageRecordList := user.HasPermission(identity.GetConponentName(), "GET", "/gui/repository/imagerecord/list")
	hasGuiRepositoryImageInformationUpgrade := user.HasPermission(identity.GetConponentName(), "GET", "/gui/repository/imageinformation/upgrade")
	hasGuiRepositoryImageInformationLog := user.HasPermission(identity.GetConponentName(), "GET", "/gui/repository/imageinformation/log")
	hasGuiDeployDeployCreate := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/deploy/create")
	hasGuiDeployDeployBlueGreenSelect := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/deploybluegreen/select")
	hasGuiRepositoryImageInformationDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/repository/imageinformation/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imageinformations/"

	imageInformationSlice := make([]ImageInformation, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &imageInformationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(imageInformationSlice); i++ {
			if hasGuiRepositoryImageRecordList {
				imageInformationSlice[i].HiddenTagGuiRepositoryImageRecordList = "<div class='btn-group'>"
			} else {
				imageInformationSlice[i].HiddenTagGuiRepositoryImageRecordList = "<div hidden>"
			}
			if hasGuiRepositoryImageInformationUpgrade {
				imageInformationSlice[i].HiddenTagGuiRepositoryImageInformationUpgrade = "<div class='btn-group'>"
			} else {
				imageInformationSlice[i].HiddenTagGuiRepositoryImageInformationUpgrade = "<div hidden>"
			}
			if hasGuiRepositoryImageInformationLog {
				imageInformationSlice[i].HiddenTagGuiRepositoryImageInformationLog = "<div class='btn-group'>"
			} else {
				imageInformationSlice[i].HiddenTagGuiRepositoryImageInformationLog = "<div hidden>"
			}
			if hasGuiDeployDeployCreate {
				imageInformationSlice[i].HiddenTagGuiDeployDeployCreate = "<div class='btn-group'>"
			} else {
				imageInformationSlice[i].HiddenTagGuiDeployDeployCreate = "<div hidden>"
			}
			if hasGuiDeployDeployBlueGreenSelect {
				imageInformationSlice[i].HiddenTagGuiDeployDeployBlueGreenSelect = "<div class='btn-group'>"
			} else {
				imageInformationSlice[i].HiddenTagGuiDeployDeployBlueGreenSelect = "<div hidden>"
			}
			if hasGuiRepositoryImageInformationDelete {
				imageInformationSlice[i].HiddenTagGuiRepositoryImageInformationDelete = "<div class='btn-group'>"
			} else {
				imageInformationSlice[i].HiddenTagGuiRepositoryImageInformationDelete = "<div hidden>"
			}
		}

		sort.Sort(ByImageInformation(imageInformationSlice))
		c.Data["imageInformationSlice"] = imageInformationSlice
	}

	guimessage.OutputMessage(c.Data)
}
