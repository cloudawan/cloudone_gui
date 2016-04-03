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
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_utility/restclient"
)

type EditController struct {
	beego.Controller
}

// @Title get
// @Description get the glusterfs cluster configuration
// @Param clusterName path string true "The name of glusterfs cluster configuration"
// @Success 200 {object} guirestapi.filesystem.glusterfs.cluster.GlusterfsCluster
// @Failure 404 error reason
// @router /:clusterName [get]
func (c *EditController) Get() {
	clusterName := c.GetString(":clusterName")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfs/clusters/" + clusterName

	glusterfsCluster := GlusterfsCluster{}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &glusterfsCluster, tokenHeaderMap)

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
	} else {
		c.Data["json"] = glusterfsCluster
		c.ServeJSON()
	}
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

// @Title create
// @Description create gluster cluster configuration
// @Param body body guirestapi.filesystem.glusterfs.cluster.GlusterfsClusterInput true "body for glusterfs cluster configuration"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [post]
func (c *EditController) Post() {
	inputBody := c.Ctx.Input.RequestBody
	glusterfsClusterInput := GlusterfsClusterInput{}
	err := json.Unmarshal(inputBody, &glusterfsClusterInput)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfs/clusters/"

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPostWithStructure(url, glusterfsClusterInput, nil, tokenHeaderMap)

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
	} else {
		c.Data["json"] = make(map[string]interface{})
		c.ServeJSON()
	}
}

// @Title update
// @Description update gluster cluster configuration
// @Param body body guirestapi.filesystem.glusterfs.cluster.GlusterfsClusterInput true "body for glusterfs cluster configuration"
// @Success 200 {string} {}
// @Failure 404 error reason
// @router / [put]
func (c *EditController) Put() {
	inputBody := c.Ctx.Input.RequestBody
	glusterfsClusterInput := GlusterfsClusterInput{}
	err := json.Unmarshal(inputBody, &glusterfsClusterInput)
	if err != nil {
		// Error
		c.Data["json"] = make(map[string]interface{})
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.Ctx.Output.Status = 404
		c.ServeJSON()
		return
	}

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/glusterfs/clusters/" + glusterfsClusterInput.Name

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPutWithStructure(url, glusterfsClusterInput, nil, tokenHeaderMap)

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
	} else {
		c.Data["json"] = make(map[string]interface{})
		c.ServeJSON()
	}
}
