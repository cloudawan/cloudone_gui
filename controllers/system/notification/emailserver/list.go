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

package emailserver

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

type EmailServerSMTP struct {
	Name                                            string
	Account                                         string
	Password                                        string
	Host                                            string
	Port                                            int
	HiddenTagGuiSystemNotificationEmailserverDelete string
}

type ByEmailServerSMTP []EmailServerSMTP

func (b ByEmailServerSMTP) Len() int           { return len(b) }
func (b ByEmailServerSMTP) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByEmailServerSMTP) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "system/notification/emailserver/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// System Notification tab menu
	user, _ := c.GetSession("user").(*rbac.User)
	c.Data["systemNotificationTabMenu"] = identity.GetSystemNotificationTabMenu(user, "emailserver")
	// Authorization for Button
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiSystemNotificationEmailServerCreate", user, "GET", "/gui/system/notification/emailserver/create")
	// Tag won't work in loop so need to be placed in data
	hasGuiSystemNotificationEmailserverDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/notification/emailserver/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/emailserversmtp/"

	emailServerSMTPSlice := make([]EmailServerSMTP, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &emailServerSMTPSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		for i := 0; i < len(emailServerSMTPSlice); i++ {
			if hasGuiSystemNotificationEmailserverDelete {
				emailServerSMTPSlice[i].HiddenTagGuiSystemNotificationEmailserverDelete = "<div class='btn-group'>"
			} else {
				emailServerSMTPSlice[i].HiddenTagGuiSystemNotificationEmailserverDelete = "<div hidden>"
			}
		}

		sort.Sort(ByEmailServerSMTP(emailServerSMTPSlice))
		c.Data["emailServerSMTPSlice"] = emailServerSMTPSlice
	}

	guimessage.OutputMessage(c.Data)
}
