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

package thirdparty

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type ListController struct {
	beego.Controller
}

type ThirdPartyApplication struct {
	Name        string
	Description string
}

func (c *ListController) Get() {
	c.TplNames = "repository/thirdparty/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/clusterapplications/"

	thirdPartyApplicationSlice := make([]ThirdPartyApplication, 0)

	_, err := restclient.RequestGetWithStructure(url, &thirdPartyApplicationSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		c.Data["thirdPartyApplicationSlice"] = thirdPartyApplicationSlice
	}

	guimessage.OutputMessage(c.Data)
}
