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
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"time"
)

type ListController struct {
	beego.Controller
}

type KubernetesEvent struct {
	Namespace      string
	Name           string
	Kind           string
	Source         map[string]interface{}
	Id             string
	FirstTimestamp time.Time
	LastTimestamp  time.Time
	Count          int
	Message        string
	Reason         string
	Acknowledge    bool
}

const (
	amountPerPage = 10
)

// @Title get
// @Description get all events and related parameters
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [get]
func (c *ListController) GetAll() {
	cloudoneAnalysisProtocol := beego.AppConfig.String("cloudoneAnalysisProtocol")
	cloudoneAnalysisHost := beego.AppConfig.String("cloudoneAnalysisHost")
	cloudoneAnalysisPort := beego.AppConfig.String("cloudoneAnalysisPort")

	acknowledge := c.GetString("acknowledge")
	if acknowledge == "" {
		acknowledge = "false"
	}

	offset, _ := c.GetInt("offset")

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
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
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
		if err == nil {
			firstTime = firstTime.Local()
		}
		lastTime, err := time.Parse(time.RFC3339, lastTimestamp)
		if err == nil {
			lastTime = lastTime.Local()
		}

		kubernetesEvent := KubernetesEvent{
			namespace,
			name,
			kind,
			source,
			id,
			firstTime,
			lastTime,
			int(count),
			message,
			reason,
			acknowledge,
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

	nextFrom := nextOffset + 1
	nextTo := nextOffset + amountPerPage

	c.Data["json"] = make(map[string]interface{})
	c.Data["json"].(map[string]interface{})["offset"] = offset
	c.Data["json"].(map[string]interface{})["amountPerPage"] = amountPerPage
	c.Data["json"].(map[string]interface{})["previousFrom"] = previousFrom
	c.Data["json"].(map[string]interface{})["previousTo"] = previousTo
	c.Data["json"].(map[string]interface{})["nextFrom"] = nextFrom
	c.Data["json"].(map[string]interface{})["nextTo"] = nextTo
	c.Data["json"].(map[string]interface{})["previousOffset"] = previousOffset
	c.Data["json"].(map[string]interface{})["nextOffset"] = nextOffset

	c.Data["json"].(map[string]interface{})["kubernetesEventSlice"] = kubernetesEventSlice
	c.ServeJSON()
}
