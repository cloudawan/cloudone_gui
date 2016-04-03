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

package imagerecord

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
)

type ListController struct {
	beego.Controller
}

type ImageRecord struct {
	ImageInformation string
	Version          string
	Path             string
	VersionInfo      map[string]string
	Environment      map[string]string
	Description      string
	CreatedTime      string
}

type ByImageRecord []ImageRecord

func (b ByImageRecord) Len() int           { return len(b) }
func (b ByImageRecord) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByImageRecord) Less(i, j int) bool { return b[i].Version > b[j].Version } // Use > to list from latest to oldest

func (c *ListController) Get() {
	c.TplName = "repository/imagerecord/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imagerecords/" + name

	imageRecordSlice := make([]ImageRecord, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &imageRecordSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		sort.Sort(ByImageRecord(imageRecordSlice))
		c.Data["imageRecordSlice"] = imageRecordSlice
	}

	guimessage.OutputMessage(c.Data)
}
