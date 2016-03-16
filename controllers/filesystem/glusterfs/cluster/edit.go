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

package cluster

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strings"
	"time"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplName = "filesystem/glusterfs/cluster/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	clusterName := c.GetString("clusterName")

	if clusterName == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Cluster Configuration"
		c.Data["createOrUpdate"] = "create"

		c.Data["sSHDialTimeoutInMilliSecond"] = 1000
		c.Data["sSHSessionTimeoutInMilliSecond"] = 10000
		c.Data["sSHPort"] = 22
	} else {
		cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
		cloudoneHost := beego.AppConfig.String("cloudoneHost")
		cloudonePort := beego.AppConfig.String("cloudonePort")

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/glusterfs/clusters/" + clusterName

		glusterfsCluster := GlusterfsCluster{}
		_, err := restclient.RequestGetWithStructure(url, &glusterfsCluster)

		if err != nil {
			// Error
			guimessage.AddDanger("Fail to get the cluster date with error " + err.Error())
			guimessage.OutputMessage(c.Data)
			return
		}

		hostList := ""
		length := len(glusterfsCluster.HostSlice)
		for index, host := range glusterfsCluster.HostSlice {
			if index == length-1 {
				hostList += host
			} else {
				hostList += host + ","
			}
		}

		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Cluster Configuration"
		c.Data["createOrUpdate"] = "update"

		c.Data["name"] = glusterfsCluster.Name
		c.Data["hostList"] = hostList
		c.Data["path"] = glusterfsCluster.Path
		c.Data["sSHDialTimeoutInMilliSecond"] = int64(glusterfsCluster.SSHDialTimeout / time.Millisecond)
		c.Data["sSHSessionTimeoutInMilliSecond"] = int64(glusterfsCluster.SSHSessionTimeout / time.Millisecond)
		c.Data["sSHPort"] = glusterfsCluster.SSHPort
		c.Data["sSHUser"] = glusterfsCluster.SSHUser
		c.Data["sSHPassword"] = glusterfsCluster.SSHPassword

		c.Data["nameFieldReadOnly"] = "readonly"
	}

	guimessage.OutputMessage(c.Data)
}

type GlusterfsClusterInput struct {
	Name                           string
	HostSlice                      []string
	Path                           string
	SSHDialTimeoutInMilliSecond    int
	SSHSessionTimeoutInMilliSecond int
	SSHPort                        int
	SSHUser                        string
	SSHPassword                    string
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")
	hostList := c.GetString("hostList")
	path := c.GetString("path")
	sSHDialTimeoutInMilliSecond, _ := c.GetInt("sSHDialTimeoutInMilliSecond")
	sSHSessionTimeoutInMilliSecond, _ := c.GetInt("sSHSessionTimeoutInMilliSecond")
	sSHPort, _ := c.GetInt("sSHPort")
	sSHUser := c.GetString("sSHUser")
	sSHPassword := c.GetString("sSHPassword")
	createOrUpdate := c.GetString("createOrUpdate")

	hostSlice := make([]string, 0)
	splitSlice := strings.Split(hostList, ",")
	for _, split := range splitSlice {
		host := strings.TrimSpace(split)
		if len(host) > 0 {
			hostSlice = append(hostSlice, host)
		}
	}

	glusterfsClusterInput := GlusterfsClusterInput{
		name,
		hostSlice,
		path,
		sSHDialTimeoutInMilliSecond,
		sSHSessionTimeoutInMilliSecond,
		sSHPort,
		sSHUser,
		sSHPassword}

	if createOrUpdate == "create" {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/glusterfs/clusters/"

		_, err := restclient.RequestPostWithStructure(url, glusterfsClusterInput, nil)

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
		} else {
			guimessage.AddSuccess("Glusterfs cluster " + name + " is created")
		}
	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/glusterfs/clusters/" + name

		_, err := restclient.RequestPutWithStructure(url, glusterfsClusterInput, nil)

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
		} else {
			guimessage.AddSuccess("Glusterfs cluster " + name + " is updated")
		}
	}

	c.Ctx.Redirect(302, "/gui/filesystem/glusterfs/cluster/")

	guimessage.RedirectMessage(c)
}
