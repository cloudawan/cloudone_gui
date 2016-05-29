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

package volume

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/filesystem/glusterfs/cluster"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strings"
)

type GlusterfsVolumeCreateParameter struct {
	ClusterName  string
	VolumeName   string
	Stripe       int
	Replica      int
	Arbiter      int
	Disperse     int
	DisperseData int
	Redundancy   int
	Transport    string
	HostSlice    []string
}

type CreateController struct {
	beego.Controller
}

func (c *CreateController) Get() {
	c.TplName = "filesystem/glusterfs/volume/create.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	clusterName := c.GetString("clusterName")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfs/clusters/" + clusterName

	glusterfsCluster := cluster.GlusterfsCluster{}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &glusterfsCluster, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		hostList := ""
		length := len(glusterfsCluster.HostSlice)
		for index, host := range glusterfsCluster.HostSlice {
			if index == length-1 {
				hostList += host
			} else {
				hostList += host + ","
			}
		}

		c.Data["hostSlice"] = glusterfsCluster.HostSlice
		c.Data["clusterName"] = clusterName
		c.Data["hostList"] = hostList
	}

	guimessage.OutputMessage(c.Data)
}

func (c *CreateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	clusterName := c.GetString("clusterName")
	hostList := c.GetString("hostList")

	name := c.GetString("name")
	stripe, _ := c.GetInt("stripe")
	replica, _ := c.GetInt("replica")
	arbiter, _ := c.GetInt("arbiter")
	disperse, _ := c.GetInt("disperse")
	disperseData, _ := c.GetInt("disperseData")
	redundancy, _ := c.GetInt("redundancy")
	transport := c.GetString("transport")

	allHostSlice := strings.Split(hostList, ",")

	hostSlice := make([]string, 0)
	for _, host := range allHostSlice {
		hostSelected := c.GetString(host)
		if hostSelected == "on" {
			hostSlice = append(hostSlice, host)
		}
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfs/clusters/" + clusterName + "/volumes/"

	glusterfsVolumeCreateParameter := GlusterfsVolumeCreateParameter{
		clusterName,
		name,
		stripe,
		replica,
		arbiter,
		disperse,
		disperseData,
		redundancy,
		transport,
		hostSlice,
	}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestPostWithStructure(url, glusterfsVolumeCreateParameter, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		guimessage.AddSuccess("Glusterfs volume " + name + " is created and started")
	}

	c.Ctx.Redirect(302, "/gui/filesystem/glusterfs/volume/list?clusterName="+clusterName)

	guimessage.RedirectMessage(c)
}
