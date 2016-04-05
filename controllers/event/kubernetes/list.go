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

package kubernetes

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

type KubernetesEvent struct {
	Namespace                              string
	Name                                   string
	Kind                                   string
	Source                                 map[string]interface{}
	Id                                     string
	FirstTimestamp                         string
	LastTimestamp                          string
	Count                                  int
	Message                                string
	Reason                                 string
	Acknowledge                            bool
	Action                                 string
	Button                                 string
	HiddenTagGuiEventKubernetesAcknowledge string
}

const (
	amountPerPage = 10
)

func (c *ListController) Get() {
	c.TplName = "event/kubernetes/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Authorization for Button
	user, _ := c.GetSession("user").(*rbac.User)
	// Tag won't work in loop so need to be placed in data
	hasGuiEventKubernetesAcknowledge := user.HasPermission(identity.GetConponentName(), "GET", "/gui/event/kubernetes/acknowledge")

	cloudoneAnalysisProtocol := beego.AppConfig.String("cloudoneAnalysisProtocol")
	cloudoneAnalysisHost := beego.AppConfig.String("cloudoneAnalysisHost")
	cloudoneAnalysisPort := beego.AppConfig.String("cloudoneAnalysisPort")

	acknowledge := c.GetString("acknowledge")
	if acknowledge == "" {
		acknowledge = "false"
	}

	offset, _ := c.GetInt("offset")

	timeZoneOffset, _ := c.GetSession("timeZoneOffset").(int)

	url := cloudoneAnalysisProtocol + "://" + cloudoneAnalysisHost + ":" + cloudoneAnalysisPort +
		"/api/v1/historicalevents?acknowledge=" + acknowledge + "&size=" + strconv.Itoa(amountPerPage) + "&offset=" + strconv.Itoa(offset)

	jsonMapSlice := make([]map[string]interface{}, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &jsonMapSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		var action string
		var button string
		if acknowledge == "true" {
			action = "false"
			button = "Unacknowledge"
		} else {
			action = "true"
			button = "Acknowledge"
		}

		kubernetesEventSlice := make([]KubernetesEvent, 0)
		for _, jsonMap := range jsonMapSlice {
			sourceJsonMap, _ := jsonMap["_source"].(map[string]interface{})

			namespace, _ := sourceJsonMap["metadata"].(map[string]interface{})["namespace"].(string)
			name, _ := sourceJsonMap["involvedObject"].(map[string]interface{})["name"].(string)
			kind, _ := sourceJsonMap["involvedObject"].(map[string]interface{})["kind"].(string)
			source, _ := sourceJsonMap["source"].(map[string]interface{})
			id, _ := jsonMap["_id"].(string)
			firstTimestamp, _ := sourceJsonMap["firstTimestamp"].(string)
			lastTimestamp, _ := sourceJsonMap["lastTimestamp"].(string)
			count, _ := sourceJsonMap["count"].(float64)
			message, _ := sourceJsonMap["message"].(string)
			reason, _ := sourceJsonMap["reason"].(string)
			acknowledge, _ := sourceJsonMap["searchMetaData"].(map[string]interface{})["acknowledge"].(bool)

			firstTime, err := time.Parse(time.RFC3339, firstTimestamp)
			if err != nil {
				// Fail to parse, show original one
			} else {
				firstTimestamp = firstTime.Add(time.Minute * time.Duration(timeZoneOffset) * -1).Format(time.RFC3339)
			}

			lastTime, err := time.Parse(time.RFC3339, lastTimestamp)
			if err != nil {
				// Fail to parse, show original one
			} else {
				lastTimestamp = lastTime.Add(time.Minute * time.Duration(timeZoneOffset) * -1).Format(time.RFC3339)
			}

			kubernetesEvent := KubernetesEvent{
				namespace,
				name,
				kind,
				source,
				id,
				firstTimestamp,
				lastTimestamp,
				int(count),
				message,
				reason,
				acknowledge,
				action,
				button,
				"",
			}

			kubernetesEventSlice = append(kubernetesEventSlice, kubernetesEvent)
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

		if acknowledge == "true" {
			c.Data["acknowledgeActive"] = "active"
			c.Data["paginationUrlPrevious"] = "/gui/event/kubernetes/list?acknowledge=true&offset=" + strconv.Itoa(previousOffset)
			c.Data["paginationUrlNext"] = "/gui/event/kubernetes/list?acknowledge=true&offset=" + strconv.Itoa(nextOffset)
		} else {
			c.Data["unacknowledgeActive"] = "active"
			c.Data["paginationUrlPrevious"] = "/gui/event/kubernetes/list?acknowledge=false&offset=" + strconv.Itoa(previousOffset)
			c.Data["paginationUrlNext"] = "/gui/event/kubernetes/list?acknowledge=false&offset=" + strconv.Itoa(nextOffset)
		}

		for i := 0; i < len(kubernetesEventSlice); i++ {
			if hasGuiEventKubernetesAcknowledge {
				kubernetesEventSlice[i].HiddenTagGuiEventKubernetesAcknowledge = "<div class='btn-group'>"
			} else {
				kubernetesEventSlice[i].HiddenTagGuiEventKubernetesAcknowledge = "<div hidden>"
			}
		}

		c.Data["kubernetesEventSlice"] = kubernetesEventSlice
	}

	guimessage.OutputMessage(c.Data)
}
