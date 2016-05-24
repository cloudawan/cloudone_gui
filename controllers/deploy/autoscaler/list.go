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

package autoscaler

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"time"
)

type ListController struct {
	beego.Controller
}

type ReplicationControllerAutoScaler struct {
	Check                              bool
	CoolDownDuration                   time.Duration
	RemainingCoolDown                  time.Duration
	KubeApiServerEndPoint              string
	KubeApiServerToken                 string
	Namespace                          string
	Kind                               string
	Name                               string
	MaximumReplica                     int
	MinimumReplica                     int
	IndicatorSlice                     []Indicator
	HiddenTagGuiDeployAutoScalerEdit   string
	HiddenTagGuiDeployAutoScalerDelete string
}

type Indicator struct {
	Type                  string
	AboveAllOrOne         bool
	AbovePercentageOfData float64
	AboveThreshold        int64
	BelowAllOrOne         bool
	BelowPercentageOfData float64
	BelowThreshold        int64
}

type ByReplicationControllerAutoScaler []ReplicationControllerAutoScaler

func (b ByReplicationControllerAutoScaler) Len() int           { return len(b) }
func (b ByReplicationControllerAutoScaler) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByReplicationControllerAutoScaler) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "deploy/autoscaler/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiDeployAutoScalerEdit", user, "GET", "/gui/deploy/autoscaler/edit")
	// Tag won't work in loop so need to be placed in data
	hasGuiDeployAutoScalerEdit := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/autoscaler/edit")
	hasGuiDeployAutoScalerDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/autoscaler/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/autoscalers/"

	replicationControllerAutoScalerSlice := make([]ReplicationControllerAutoScaler, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &replicationControllerAutoScalerSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(replicationControllerAutoScalerSlice); i++ {
			if hasGuiDeployAutoScalerEdit {
				replicationControllerAutoScalerSlice[i].HiddenTagGuiDeployAutoScalerEdit = "<div class='btn-group'>"
			} else {
				replicationControllerAutoScalerSlice[i].HiddenTagGuiDeployAutoScalerEdit = "<div hidden>"
			}
			if hasGuiDeployAutoScalerDelete {
				replicationControllerAutoScalerSlice[i].HiddenTagGuiDeployAutoScalerDelete = "<div class='btn-group'>"
			} else {
				replicationControllerAutoScalerSlice[i].HiddenTagGuiDeployAutoScalerDelete = "<div hidden>"
			}
		}

		sort.Sort(ByReplicationControllerAutoScaler(replicationControllerAutoScalerSlice))
		c.Data["replicationControllerAutoScalerSlice"] = replicationControllerAutoScalerSlice
	}

	guimessage.OutputMessage(c.Data)
}
