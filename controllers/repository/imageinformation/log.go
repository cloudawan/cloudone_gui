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

package imageinformation

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
)

type LogController struct {
	beego.Controller
}

func (c *LogController) Get() {
	c.TplName = "repository/imageinformation/log.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	logKey := c.GetString("logKey")

	if logKey == "" {
		guimessage.AddDanger("No log key is passed")
		guimessage.OutputMessage(c.Data)
		return
	}

	outputMessage := c.GetSession(logKey)

	c.DelSession(logKey)

	c.Data["log"] = outputMessage

	guimessage.OutputMessage(c.Data)
}
