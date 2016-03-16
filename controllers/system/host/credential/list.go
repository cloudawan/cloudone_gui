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

package credential

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
)

type ListController struct {
	beego.Controller
}

type Credential struct {
	IP  string
	SSH SSH
}

type SSH struct {
	Port     int
	User     string
	Password string
}

type ByCredential []Credential

func (b ByCredential) Len() int           { return len(b) }
func (b ByCredential) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByCredential) Less(i, j int) bool { return b[i].IP < b[j].IP }

func (c *ListController) Get() {
	c.TplName = "system/host/credential/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/hosts/credentials/"

	credentialSlice := make([]Credential, 0)

	_, err := restclient.RequestGetWithStructure(url, &credentialSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		sort.Sort(ByCredential(credentialSlice))
		c.Data["credentialSlice"] = credentialSlice
	}

	guimessage.OutputMessage(c.Data)
}
