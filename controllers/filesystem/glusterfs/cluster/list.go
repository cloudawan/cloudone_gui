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

package cluster

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"time"
)

type ListController struct {
	beego.Controller
}

type GlusterfsCluster struct {
	Name              string
	HostSlice         []string
	Path              string
	SSHDialTimeout    time.Duration
	SSHSessionTimeout time.Duration
	SSHPort           int
	SSHUser           string
	SSHPassword       string
}

type ByGlusterfsCluster []GlusterfsCluster

func (b ByGlusterfsCluster) Len() int           { return len(b) }
func (b ByGlusterfsCluster) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByGlusterfsCluster) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "filesystem/glusterfs/cluster/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfs/clusters/"

	glusterfsClusterSlice := make([]GlusterfsCluster, 0)

	_, err := restclient.RequestGetWithStructure(url, &glusterfsClusterSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		sort.Sort(ByGlusterfsCluster(glusterfsClusterSlice))
		c.Data["glusterfsClusterSlice"] = glusterfsClusterSlice
	}

	guimessage.OutputMessage(c.Data)
}
