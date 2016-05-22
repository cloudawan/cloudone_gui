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
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"math"
	"strconv"
	"strings"
)

const (
	componentName = "cloudone_gui"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "identity/login.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	guimessage.OutputMessage(c.Data)
}

type UserData struct {
	Username string
	Password string
}

type TokenData struct {
	Token string
}

func (c *LoginController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	username := c.GetString("username")
	password := c.GetString("password")
	timeZoneOffset, err := c.GetInt("timeZoneOffset")
	if err != nil {
		guimessage.AddDanger("Fail to get browser time zone offset. Use UTC instead")
	} else {
		hourOffset := float64(timeZoneOffset) / 60.0
		// Since it is time offset, it needs to multiple -1 to get the UTC format
		sign := "-"
		if hourOffset < 0 {
			sign = "+"
		}
		guimessage.AddSuccess("Browser time zone is UTC " + sign + strconv.FormatFloat(math.Abs(hourOffset), 'f', -1, 64) + "")
		c.SetSession("timeZoneOffset", timeZoneOffset)
	}

	// User
	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/authorizations/tokens/"
	userData := UserData{username, password}
	tokenData := TokenData{}
	errorInterface, err := restclient.RequestPostWithStructure(url, userData, &tokenData, nil)
	errorJsonMap, _ := errorInterface.(map[string]interface{})

	if err != nil {
		if errorJsonMap != nil {
			errorMessage, _ := errorJsonMap["ErrorMessage"].(string)
			guimessage.AddDanger(errorMessage)
		} else {
			guimessage.AddDanger(err.Error())
		}

		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/login/")
		return
	}

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/authorizations/tokens/" + tokenData.Token + "/components/" + componentName
	user := &rbac.User{}
	_, err = restclient.RequestGetWithStructure(url, &user, nil)
	if err != nil {
		guimessage.AddDanger("Fail to get user data with error " + err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/login/")
		return
	}

	// Set session
	// User is used by this component to authorize
	c.SetSession("user", user)
	c.SetSession("username", user.Name)
	// Token is used to submit to other componentes to authorize
	headerMap := make(map[string]string)
	headerMap["token"] = tokenData.Token
	c.SetSession("tokenHeaderMap", headerMap)
	// Layout menu is used to display common layout menu
	layoutMenu := GetLayoutMenu(user)
	c.SetSession("layoutMenu", layoutMenu)

	// Namespace
	namespace := beego.AppConfig.String("namespace")
	// Get first namespace
	for _, resource := range user.ResourceSlice {
		if resource.Component == componentName || resource.Component == "*" {
			if strings.HasPrefix(resource.Path, "/namespaces/") {
				splitSlice := strings.Split(resource.Path, "/")
				name := splitSlice[2]
				if name != "" && name != "*" {
					namespace = name
					break
				}
			}
		}
	}
	metaDataMap := user.MetaDataMap
	if metaDataMap == nil {
		metaDataMap = make(map[string]string)
	}
	// If loginNamespace is set, use it
	loginNamespace := metaDataMap["loginNamespace"]
	if len(loginNamespace) > 0 {
		namespace = loginNamespace
	}

	// Check if the namespace existing
	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort + "/api/v1/namespaces/"

	nameSlice := make([]string, 0)

	_, err = restclient.RequestGetWithStructure(url, &nameSlice, headerMap)

	if IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		guimessage.AddDanger("Fail to get namespace data with error " + err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/login/")
		return
	}

	namespaceExist := false
	for _, name := range nameSlice {
		if name == namespace {
			namespaceExist = true
		}
	}

	if namespaceExist == false {
		guimessage.AddDanger("Namespace " + namespace + " set for user doesn't exist")
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/login/")
		return
	}

	// Set namespace
	c.SetSession("namespace", namespace)

	// Send audit log since this page will pass filter
	go func() {
		sendAuditLog(c.Ctx, user.Name, false)
	}()

	guimessage.AddSuccess("User " + username + " login")

	c.Ctx.Redirect(302, "/gui/dashboard/topology/")

	guimessage.RedirectMessage(c)
}
