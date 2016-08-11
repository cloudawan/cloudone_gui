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

package daemon

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

type SLBDaemon struct {
	Name                                 string
	EndPointSlice                        []string
	NodeHostSlice                        []string
	Description                          string
	HiddenTagGuiSystemSLBDaemonEdit      string
	HiddenTagGuiSystemSLBDaemonConfigure string
	HiddenTagGuiSystemSLBDaemonDelete    string
}

type BySLBDaemon []SLBDaemon

func (b BySLBDaemon) Len() int           { return len(b) }
func (b BySLBDaemon) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b BySLBDaemon) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "system/slb/daemon/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiSystemSLBDaemonEdit", user, "GET", "/gui/system/slb/daemon/edit")
	// Tag won't work in loop so need to be placed in data
	hasGuiSystemSLBDaemonEdit := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/slb/daemon/edit")
	hasGuiSystemSLBDaemonConfigure := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/slb/daemon/configure")
	hasGuiSystemSLBDaemonDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/system/slb/daemon/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/slbs/daemons/"

	slbDaemonSlice := make([]SLBDaemon, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &slbDaemonSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		for i := 0; i < len(slbDaemonSlice); i++ {
			if hasGuiSystemSLBDaemonEdit {
				slbDaemonSlice[i].HiddenTagGuiSystemSLBDaemonEdit = "<div class='btn-group'>"
			} else {
				slbDaemonSlice[i].HiddenTagGuiSystemSLBDaemonEdit = "<div hidden>"
			}
			if hasGuiSystemSLBDaemonConfigure {
				slbDaemonSlice[i].HiddenTagGuiSystemSLBDaemonConfigure = "<div class='btn-group'>"
			} else {
				slbDaemonSlice[i].HiddenTagGuiSystemSLBDaemonConfigure = "<div hidden>"
			}
			if hasGuiSystemSLBDaemonDelete {
				slbDaemonSlice[i].HiddenTagGuiSystemSLBDaemonDelete = "<div class='btn-group'>"
			} else {
				slbDaemonSlice[i].HiddenTagGuiSystemSLBDaemonDelete = "<div hidden>"
			}
		}

		sort.Sort(BySLBDaemon(slbDaemonSlice))
		c.Data["slbDaemonSlice"] = slbDaemonSlice
	}

	guimessage.OutputMessage(c.Data)
}
