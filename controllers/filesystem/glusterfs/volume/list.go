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
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type ListController struct {
	beego.Controller
}

type GlusterfsVolume struct {
	VolumeName                                  string
	Type                                        string
	VolumeID                                    string
	Status                                      string
	NumberOfBricks                              string
	TransportType                               string
	Bricks                                      []string
	Size                                        int
	ClusterName                                 string
	HiddenTagGuiFileSystemGlusterfsVolumeReset  string
	HiddenTagGuiFileSystemGlusterfsVolumeDelete string
}

func (c *ListController) Get() {
	c.TplName = "filesystem/glusterfs/volume/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuiFileSystemGlusterfsClusterList", user, "GET", "/gui/filesystem/glusterfs/cluster/list")
	identity.SetPrivilegeHiddenTag(c.Data, "hiddenTagGuifFileSystemGlusterfsVolumeCreate", user, "GET", "/gui/filesystem/glusterfs/volume/create")
	// Tag won't work in loop so need to be placed in data
	hasHiddenTagGuiFileSystemGlusterfsVolumeReset := user.HasPermission(identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/volume/reset")
	hasHiddenTagGuiFileSystemGlusterfsVolumeDelete := user.HasPermission(identity.GetConponentName(), "GET", "/gui/filesystem/glusterfs/volume/delete")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	clusterName := c.GetString("clusterName")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfs/clusters/" + clusterName + "/volumes/"

	glusterfsVolumeSlice := make([]GlusterfsVolume, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &glusterfsVolumeSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		for i := 0; i < len(glusterfsVolumeSlice); i++ {
			glusterfsVolumeSlice[i].ClusterName = clusterName

			if hasHiddenTagGuiFileSystemGlusterfsVolumeReset {
				glusterfsVolumeSlice[i].HiddenTagGuiFileSystemGlusterfsVolumeReset = "<div class='btn-group'>"
			} else {
				glusterfsVolumeSlice[i].HiddenTagGuiFileSystemGlusterfsVolumeReset = "<div hidden>"
			}
			if hasHiddenTagGuiFileSystemGlusterfsVolumeDelete {
				glusterfsVolumeSlice[i].HiddenTagGuiFileSystemGlusterfsVolumeDelete = "<div class='btn-group'>"
			} else {
				glusterfsVolumeSlice[i].HiddenTagGuiFileSystemGlusterfsVolumeDelete = "<div hidden>"
			}
		}

		c.Data["glusterfsVolumeSlice"] = glusterfsVolumeSlice
		c.Data["clusterName"] = clusterName
	}

	guimessage.OutputMessage(c.Data)
}
