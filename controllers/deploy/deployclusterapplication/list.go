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

package deployclusterapplication

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

type DeployClusterApplication struct {
	Name                                             string
	Size                                             int
	ServiceName                                      string
	ReplicationControllerNameSlice                   []string
	HiddenTagGuiDeployDeployClusterApplicationSize   string
	HiddenTagGuiDeployDeployClusterApplicationDelete string
}

type ByDeployClusterApplication []DeployClusterApplication

func (b ByDeployClusterApplication) Len() int           { return len(b) }
func (b ByDeployClusterApplication) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByDeployClusterApplication) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "deploy/deployclusterapplication/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	// Tag won't work in loop so need to be placed in data
	hasGuiDeployDeployClusterApplicationSize := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/deployclusterapplication/size")
	hasGuiDeployDeployClusterApplicationDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/deploy/deployclusterapplication/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespace := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deployclusterapplications/" + namespace

	deployClusterApplicationSlice := make([]DeployClusterApplication, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &deployClusterApplicationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(deployClusterApplicationSlice); i++ {
			if hasGuiDeployDeployClusterApplicationSize {
				deployClusterApplicationSlice[i].HiddenTagGuiDeployDeployClusterApplicationSize = "<div class='btn-group'>"
			} else {
				deployClusterApplicationSlice[i].HiddenTagGuiDeployDeployClusterApplicationSize = "<div hidden>"
			}
			if hasGuiDeployDeployClusterApplicationDelete {
				deployClusterApplicationSlice[i].HiddenTagGuiDeployDeployClusterApplicationDelete = "<div class='btn-group'>"
			} else {
				deployClusterApplicationSlice[i].HiddenTagGuiDeployDeployClusterApplicationDelete = "<div hidden>"
			}
		}

		sort.Sort(ByDeployClusterApplication(deployClusterApplicationSlice))
		c.Data["deployClusterApplicationSlice"] = deployClusterApplicationSlice
	}

	guimessage.OutputMessage(c.Data)
}
