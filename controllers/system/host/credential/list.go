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

package credential

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

type Credential struct {
	IP                                     string
	SSH                                    SSH
	HiddenTagGuiSystemHostCredentialEdit   string
	HiddenTagGuiSystemHostCredentialDelete string
}

type SSH struct {
	Port     int
	User     string
	Password string
}

type ByCredential []Credential

func (b ByCredential) Len() int           { return len(b) }
func (b ByCredential) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByCredential) Less(i, j int) bool { return b[i].IP < b[j].IP }

func (c *ListController) Get() {
	c.TplName = "system/host/credential/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPriviledgeHiddenTag(c.Data, "hiddenTagGuiSystemHostCredentialEdit", user, "GET", "/gui/system/host/credential/edit")
	// Tag won't work in loop so need to be placed in data
	hasGuiSystemHostCredentialEdit := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/host/credential/edit")
	hasGuiSystemHostCredentialDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/host/credential/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/hosts/credentials/"

	credentialSlice := make([]Credential, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &credentialSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(credentialSlice); i++ {
			if hasGuiSystemHostCredentialEdit {
				credentialSlice[i].HiddenTagGuiSystemHostCredentialEdit = "<div class='btn-group'>"
			} else {
				credentialSlice[i].HiddenTagGuiSystemHostCredentialEdit = "<div hidden>"
			}
			if hasGuiSystemHostCredentialDelete {
				credentialSlice[i].HiddenTagGuiSystemHostCredentialDelete = "<div class='btn-group'>"
			} else {
				credentialSlice[i].HiddenTagGuiSystemHostCredentialDelete = "<div hidden>"
			}
		}

		sort.Sort(ByCredential(credentialSlice))
		c.Data["credentialSlice"] = credentialSlice
	}

	guimessage.OutputMessage(c.Data)
}
