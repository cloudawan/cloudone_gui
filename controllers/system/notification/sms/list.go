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

package sms

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

type SMSNexmo struct {
	Name                                    string
	Url                                     string
	APIKey                                  string
	APISecret                               string
	HiddenTagGuiSystemNotificationSMSDelete string
}

type BySMSNexmo []SMSNexmo

func (b BySMSNexmo) Len() int           { return len(b) }
func (b BySMSNexmo) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b BySMSNexmo) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "system/notification/sms/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// System Notification tab menu
	user, _ := c.GetSession("user").(*rbac.User)
	c.Data["systemNotificationTabMenu"] = identity.GetSystemNotificationTabMenu(user, "sms")
	// Authorization for Button
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiSystemNotificationSMSCreate", user, "GET", "/gui/system/notification/sms/create")
	// Tag won't work in loop so need to be placed in data
	hasGuiSystemNotificationSMSDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/notification/sms/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/smsnexmo/"

	smsNexmoSlice := make([]SMSNexmo, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &smsNexmoSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		for i := 0; i < len(smsNexmoSlice); i++ {
			if hasGuiSystemNotificationSMSDelete {
				smsNexmoSlice[i].HiddenTagGuiSystemNotificationSMSDelete = "<div class='btn-group'>"
			} else {
				smsNexmoSlice[i].HiddenTagGuiSystemNotificationSMSDelete = "<div hidden>"
			}
		}

		sort.Sort(BySMSNexmo(smsNexmoSlice))
		c.Data["smsNexmoSlice"] = smsNexmoSlice
	}

	guimessage.OutputMessage(c.Data)
}
