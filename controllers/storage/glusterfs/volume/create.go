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
	"fmt"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type GlusterfsVolumeInput struct {
	Name      string
	Stripe    int
	Replica   int
	Transport string
	HostSlice []string
}

type GlusterfsVolumeControl struct {
	GlusterfsClusterHostSlice []string
	GlusterfsPath             string
}

type CreateController struct {
	beego.Controller
}

func (c *CreateController) Get() {
	c.TplNames = "storage/glusterfs/volume/create.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfsvolumes/configuration"

	glusterfsVolumeControl := GlusterfsVolumeControl{}
	_, err := restclient.RequestGetWithStructure(url, &glusterfsVolumeControl)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["glusterfsClusterHostSlice"] = glusterfsVolumeControl.GlusterfsClusterHostSlice
	}

	guimessage.OutputMessage(c.Data)
}

func (c *CreateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")
	stripe, _ := c.GetInt("stripe")
	replica, _ := c.GetInt("replica")
	transport := c.GetString("transport")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfsvolumes/configuration"

	glusterfsVolumeControl := GlusterfsVolumeControl{}
	_, err := restclient.RequestGetWithStructure(url, &glusterfsVolumeControl)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		hostSlice := make([]string, 0)
		for _, host := range glusterfsVolumeControl.GlusterfsClusterHostSlice {
			hostSelected := c.GetString(host)
			if hostSelected == "on" {
				hostSlice = append(hostSlice, host)
			}
		}

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/glusterfsvolumes/"

		glusterfsVolumeInput := GlusterfsVolumeInput{
			name,
			stripe,
			replica,
			transport,
			hostSlice,
		}

		fmt.Println(glusterfsVolumeInput)

		_, err := restclient.RequestPostWithStructure(url, glusterfsVolumeInput, nil)

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
		} else {
			guimessage.AddSuccess("Glusterfs volume " + name + " is created and started")
		}
	}

	c.Ctx.Redirect(302, "/gui/storage/glusterfs/volume/")

	guimessage.RedirectMessage(c)
}
