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

package identity

import (
	"github.com/astaxie/beego/context"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"strings"
)

func IsTokenInvalid(err error) bool {
	if err == nil {
		return false
	} else if strings.Contains(err.Error(), "Token is incorrect or expired") {
		return true
	} else {
		return false
	}
}

func IsTokenInvalidAndRedirect(c guimessagedisplay.SessionUtility, ctx *context.Context, err error) bool {
	if IsTokenInvalid(err) {
		guimessage := guimessagedisplay.GetGUIMessage(c)
		guimessage.AddDanger("User token is expired. Please login agin.")
		guimessage.RedirectMessage(c)

		c.DelSession("user")
		c.DelSession("tokenHeaderMap")

		ctx.Redirect(302, "/gui/login/")

		return true
	} else {
		return false
	}
}
