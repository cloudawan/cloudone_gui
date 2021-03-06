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

package namespace

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
)

type SelectController struct {
	beego.Controller
}

func (c *SelectController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	name := c.GetString("name")

	c.SetSession("namespace", name)
	guimessage.AddSuccess("Use namespace " + name)

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/system/namespace/list")

	guimessage.RedirectMessage(c)
}
