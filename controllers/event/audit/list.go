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

func (c *ListController) Get() {
	c.TplName = "event/audit/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneAnalysisProtocol := beego.AppConfig.String("cloudoneAnalysisProtocol")
	cloudoneAnalysisHost := beego.AppConfig.String("cloudoneAnalysisHost")
	cloudoneAnalysisPort := beego.AppConfig.String("cloudoneAnalysisPort")

	offset, _ := c.GetInt("offset")

	timeZoneOffset, _ := c.GetSession("timeZoneOffset").(int)

	url := cloudoneAnalysisProtocol + "://" + cloudoneAnalysisHost + ":" + cloudoneAnalysisPort +
		"/api/v1/auditlogs?size=" + strconv.Itoa(amountPerPage) + "&offset=" + strconv.Itoa(offset)

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
			auditLogSlice[i].CreatedTime.Add(time.Minute * time.Duration(timeZoneOffset) * -1)
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
		c.Data["paginationUrlPrevious"] = "/gui/event/audit/list?offset=" + strconv.Itoa(previousOffset)
		c.Data["paginationUrlNext"] = "/gui/event/audit/list?offset=" + strconv.Itoa(nextOffset)

		c.Data["auditLogSlice"] = auditLogSlice
	}

	guimessage.OutputMessage(c.Data)
}
