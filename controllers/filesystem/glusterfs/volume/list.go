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
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type ListController struct {
	beego.Controller
}

type GlusterfsVolume struct {
	VolumeName     string
	Type           string
	VolumeID       string
	Status         string
	NumberOfBricks string
	TransportType  string
	Bricks         []string
	Size           int
	ClusterName    string
}

func (c *ListController) Get() {
	c.TplNames = "filesystem/glusterfs/volume/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	clusterName := c.GetString("clusterName")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfs/clusters/" + clusterName + "/volumes/"

	glusterfsVolumeSlice := make([]GlusterfsVolume, 0)

	_, err := restclient.RequestGetWithStructure(url, &glusterfsVolumeSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(glusterfsVolumeSlice); i++ {
			glusterfsVolumeSlice[i].ClusterName = clusterName
		}

		c.Data["glusterfsVolumeSlice"] = glusterfsVolumeSlice
		c.Data["clusterName"] = clusterName
	}

	guimessage.OutputMessage(c.Data)
}
