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

package daemon

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strings"
)

type EditController struct {
	beego.Controller
}

type Region struct {
	Name           string
	LocationTagged bool
	ZoneSlice      []Zone
}

type Zone struct {
	Name           string
	LocationTagged bool
	NodeSlice      []Node
}

type Node struct {
	Name    string
	Address string
	Checked string
}

func (c *EditController) Get() {
	c.TplName = "system/slb/daemon/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	name := c.GetString("name")

	// Get node topology
	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort + "/api/v1/nodes/topology"

	regionSlice := make([]Region, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &regionSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		guimessage.OutputMessage(c.Data)
		return
	}

	if name == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create SLB Daemon"
		c.Data["createOrUpdate"] = "create"
		c.Data["regionSlice"] = regionSlice
	} else {
		cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
		cloudoneHost := beego.AppConfig.String("cloudoneHost")
		cloudonePort := beego.AppConfig.String("cloudonePort")

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/slbs/daemons/" + name

		slbDaemon := SLBDaemon{}

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestGetWithStructure(url, &slbDaemon, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
			guimessage.OutputMessage(c.Data)
			return
		}

		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update SLB Daemon"
		c.Data["createOrUpdate"] = "update"

		c.Data["name"] = slbDaemon.Name
		c.Data["description"] = slbDaemon.Description

		for _, nodeHost := range slbDaemon.NodeHostSlice {
			for i, region := range regionSlice {
				for j, zone := range region.ZoneSlice {
					for k, node := range zone.NodeSlice {
						if node.Address == nodeHost {
							regionSlice[i].ZoneSlice[j].NodeSlice[k].Checked = "checked"
						}
					}
				}
			}
		}
		c.Data["regionSlice"] = regionSlice

		endPointList := ""
		length := len(slbDaemon.EndPointSlice)
		for index, endPoint := range slbDaemon.EndPointSlice {
			if index == length-1 {
				endPointList += endPoint
			} else {
				endPointList += endPoint + ","
			}
		}
		c.Data["endPointList"] = endPointList

		c.Data["ipFieldReadOnly"] = "readonly"
	}

	guimessage.OutputMessage(c.Data)
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	createOrUpdate := c.GetString("createOrUpdate")
	name := c.GetString("name")
	description := c.GetString("description")

	nodeHostSlice := make([]string, 0)
	inputMap := c.Input()
	if inputMap != nil {
		for key, _ := range inputMap {
			// Only collect host
			if strings.HasPrefix(key, "nodeHost_") {
				nodeHostSlice = append(nodeHostSlice, key[len("nodeHost_"):])
			}
		}
	}

	endPointList := c.GetString("endPointList")
	endPointSlice := make([]string, 0)
	splitSlice := strings.Split(endPointList, ",")
	for _, split := range splitSlice {
		endPoint := strings.TrimSpace(split)
		if len(endPoint) > 0 {
			endPointSlice = append(endPointSlice, endPoint)
		}
	}

	slbDaemon := SLBDaemon{
		name,
		endPointSlice,
		nodeHostSlice,
		description,
		"",
		"",
	}

	if createOrUpdate == "create" {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/slbs/daemons/"

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestPostWithStructure(url, slbDaemon, nil, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		} else {
			guimessage.AddSuccess("SLB daemon " + name + " is created")
		}
	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/slbs/daemons/" + name

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestPutWithStructure(url, slbDaemon, nil, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		} else {
			guimessage.AddSuccess("SLB daemon " + name + " is updated")
		}
	}

	c.Ctx.Redirect(302, "/gui/system/slb/daemon/list")

	guimessage.RedirectMessage(c)
}
