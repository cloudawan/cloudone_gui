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
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type BookmarkController struct {
	beego.Controller
}

func (c *BookmarkController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")

	sessionUser, _ := c.GetSession("user").(*rbac.User)
	if sessionUser == nil {
		// Error
		guimessage.AddDanger("User in seesion does not exist")
		c.Ctx.Redirect(302, "/gui/system/namespace/list")
		guimessage.RedirectMessage(c)
		return
	}

	userName := sessionUser.Name

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/authorizations/users/" + userName

	user := rbac.User{}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &user, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		c.Ctx.Redirect(302, "/gui/system/namespace/list")
		guimessage.RedirectMessage(c)
		return
	}

	metaDataMap := user.MetaDataMap
	if metaDataMap == nil {
		metaDataMap = make(map[string]string)
	}
	metaDataMap["loginNamespace"] = name

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/authorizations/users/" + userName + "/metadata"

	_, err = restclient.RequestPutWithStructure(url, metaDataMap, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		c.Ctx.Redirect(302, "/gui/system/namespace/list")
		guimessage.RedirectMessage(c)
		return
	}

	// Set session
	sessionUser.MetaDataMap = metaDataMap
	c.SetSession("user", sessionUser)

	guimessage.AddSuccess("Bookmark the namespace " + name + " as the login namespace")

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/system/namespace/list")

	guimessage.RedirectMessage(c)
}
