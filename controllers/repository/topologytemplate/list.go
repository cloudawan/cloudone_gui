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

package topologytemplate

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"strconv"
	"time"
)

type ListController struct {
	beego.Controller
}

type Topology struct {
	Name            string
	SourceNamespace string
	CreatedUser     string
	CreatedDate     time.Time
	Description     string
	LaunchSlice     []Launch
	// The following is used for GUI
	HiddenTagGuiRepositoryTopologyTemplateClone  string
	HiddenTagGuiRepositoryTopologyTemplateDelete string
}

type Launch struct {
	Order                    int
	LaunchApplication        *LaunchApplication
	LaunchClusterApplication *LaunchClusterApplication
	// The following is used for GUI
	Information                       string
	Name                              string
	RegionSlice                       []Region
	HiddenTagLaunchApplication        string
	HiddenTagLaunchClusterApplication string
}

type ByLaunch []Launch

func (b ByLaunch) Len() int           { return len(b) }
func (b ByLaunch) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByLaunch) Less(i, j int) bool { return b[i].Order < b[j].Order }

type Region struct {
	Name           string
	LocationTagged bool
	ZoneSlice      []Zone
}

type Zone struct {
	Name           string
	LocationTagged bool
}

type LaunchApplication struct {
	ImageInformationName string
	Version              string
	Description          string
	ReplicaAmount        int
	PortSlice            []DeployContainerPort
	EnvironmentSlice     []ReplicationControllerContainerEnvironment
	ResourceMap          map[string]interface{}
	ExtraJsonMap         map[string]interface{}
}

type LaunchClusterApplication struct {
	Name                              string
	Size                              int
	EnvironmentSlice                  []interface{}
	ReplicationControllerExtraJsonMap map[string]interface{}
}

type DeployContainerPort struct {
	Name          string
	ContainerPort int
	NodePort      int
}

type ReplicationControllerContainerEnvironment struct {
	Name  string
	Value string
}

type ByTopology []Topology

func (b ByTopology) Len() int           { return len(b) }
func (b ByTopology) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByTopology) Less(i, j int) bool { return b[i].CreatedDate.Before(b[j].CreatedDate) }

func (c *ListController) Get() {
	c.TplName = "repository/topologytemplate/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	// Tag won't work in loop so need to be placed in data
	hasGuiRepositoryTopologyTemplateClone := user.HasPermission(identity.GetConponentName(), "GET", "/gui/repository/topologytemplate/clone")
	hasGuiRepositoryTopologyTemplateDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/repository/topologytemplate/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/topology/"

	topologySlice := make([]Topology, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &topologySlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(topologySlice); i++ {
			if hasGuiRepositoryTopologyTemplateClone {
				topologySlice[i].HiddenTagGuiRepositoryTopologyTemplateClone = "<div class='btn-group'>"
			} else {
				topologySlice[i].HiddenTagGuiRepositoryTopologyTemplateClone = "<div hidden>"
			}
			if hasGuiRepositoryTopologyTemplateDelete {
				topologySlice[i].HiddenTagGuiRepositoryTopologyTemplateDelete = "<div class='btn-group'>"
			} else {
				topologySlice[i].HiddenTagGuiRepositoryTopologyTemplateDelete = "<div hidden>"
			}

			for j := 0; j < len(topologySlice[i].LaunchSlice); j++ {
				if topologySlice[i].LaunchSlice[j].LaunchApplication != nil {
					topologySlice[i].LaunchSlice[j].Information = "APP: " + topologySlice[i].LaunchSlice[j].LaunchApplication.ImageInformationName + " " + topologySlice[i].LaunchSlice[j].LaunchApplication.Version + " (" + strconv.Itoa(topologySlice[i].LaunchSlice[j].LaunchApplication.ReplicaAmount) + ")"
				} else if topologySlice[i].LaunchSlice[j].LaunchClusterApplication != nil {
					topologySlice[i].LaunchSlice[j].Information = "3rd: " + topologySlice[i].LaunchSlice[j].LaunchClusterApplication.Name + " (" + strconv.Itoa(topologySlice[i].LaunchSlice[j].LaunchClusterApplication.Size) + ")"
				}

			}

			// Change time format
			topologySlice[i].CreatedDate = topologySlice[i].CreatedDate.Round(time.Second)
		}

		sort.Sort(ByTopology(topologySlice))
		c.Data["topologySlice"] = topologySlice
	}

	guimessage.OutputMessage(c.Data)
}
