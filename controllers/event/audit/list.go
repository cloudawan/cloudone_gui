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

package audit

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"time"
)

type ListController struct {
	beego.Controller
}

type AuditLog struct {
	Component         string
	Kind              string
	Path              string
	UserName          string
	RemoteAddress     string
	RemoteHost        string
	CreatedTime       time.Time
	QueryParameterMap map[string][]string
	PathParameterMap  map[string]string
	RequestMethod     string
	RequestURI        string
	RequestBody       string
	RequestHeader     map[string][]string
	Description       string
}

const (
	amountPerPage = 10
)

type UserData struct {
	Name     string
	Selected string
}

func (c *ListController) Get() {
	c.TplName = "event/audit/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneAnalysisProtocol := beego.AppConfig.String("cloudoneAnalysisProtocol")
	cloudoneAnalysisHost := beego.AppConfig.String("cloudoneAnalysisHost")
	cloudoneAnalysisPort := beego.AppConfig.String("cloudoneAnalysisPort")

	offset, _ := c.GetInt("offset")
	userName := c.GetString("userName")

	if userName == "All" {
		userName = ""
	}

	url := cloudoneAnalysisProtocol + "://" + cloudoneAnalysisHost + ":" + cloudoneAnalysisPort +
		"/api/v1/auditlogs/" + userName + "?size=" + strconv.Itoa(amountPerPage) + "&offset=" + strconv.Itoa(offset)

	auditLogSlice := make([]AuditLog, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &auditLogSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		for i := 0; i < len(auditLogSlice); i++ {
			auditLogSlice[i].CreatedTime = auditLogSlice[i].CreatedTime.In(time.Local)
		}

		previousOffset := offset - amountPerPage
		if previousOffset < 0 {
			previousOffset = 0
		}
		nextOffset := offset + amountPerPage

		previousFrom := previousOffset
		if previousFrom < 0 {
			previousFrom = 0
		}
		previousFrom += 1
		previousTo := previousOffset + amountPerPage
		c.Data["previousLabel"] = strconv.Itoa(previousFrom) + "~" + strconv.Itoa(previousTo)
		if offset == 0 {
			c.Data["previousButtonHidden"] = "hidden"
		} else {
			c.Data["previousButtonHidden"] = ""
		}

		nextFrom := nextOffset + 1
		nextTo := nextOffset + amountPerPage
		c.Data["nextLabel"] = strconv.Itoa(nextFrom) + "~" + strconv.Itoa(nextTo)

		c.Data["acknowledgeActive"] = "active"
		c.Data["paginationUrlPrevious"] = "/gui/event/audit/list?offset=" + strconv.Itoa(previousOffset) + "&userName=" + userName
		c.Data["paginationUrlNext"] = "/gui/event/audit/list?offset=" + strconv.Itoa(nextOffset) + "&userName=" + userName

		c.Data["auditLogSlice"] = auditLogSlice

		// Get user slice to select
		cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
		cloudoneHost := beego.AppConfig.String("cloudoneHost")
		cloudonePort := beego.AppConfig.String("cloudonePort")

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/authorizations/users"

		userSlice := make([]rbac.User, 0)

		_, err = restclient.RequestGetWithStructure(url, &userSlice, tokenHeaderMap)

		if err != nil {
			guimessage.AddDanger(err.Error())
		} else {
			userDataSlice := make([]UserData, 0)
			userDataSlice = append(userDataSlice, UserData{"All", ""})
			for _, user := range userSlice {
				if user.Name == userName {
					userDataSlice = append(userDataSlice, UserData{user.Name, "selected"})
				} else {
					userDataSlice = append(userDataSlice, UserData{user.Name, ""})
				}
			}
			c.Data["userDataSlice"] = userDataSlice
		}
	}

	guimessage.OutputMessage(c.Data)
}
