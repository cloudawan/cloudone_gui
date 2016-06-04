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

package imageinformation

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_gui/controllers/utility/random"
	"github.com/cloudawan/cloudone_utility/restclient"
	"regexp"
)

type CreateController struct {
	beego.Controller
}

type PrivateRegistry struct {
	Name string
	Host string
	Port int
}

func (c *CreateController) Get() {
	c.TplName = "repository/imageinformation/create.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/privateregistries/servers/"

	privateRegistrySlice := make([]PrivateRegistry, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &privateRegistrySlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/repository/imageinformation/list")
		return
	}

	if len(privateRegistrySlice) == 0 {
		guimessage.AddWarning("At least one private registry is required. Please configure it in system.")
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/repository/imageinformation/list")
		return
	}

	c.Data["privateRegistrySlice"] = privateRegistrySlice

	guimessage.OutputMessage(c.Data)
}

func (c *CreateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")
	kind := c.GetString("kind")
	description := c.GetString("description")

	// Name need to be a DNS 952 label
	match, _ := regexp.MatchString("^[a-z]{1}[a-z0-9-]{1,23}$", name)
	if match == false {
		guimessage.AddDanger("The name need to be a DNS 952 label ^[a-z]{1}[a-z0-9-]{1,23}$")
		c.Ctx.Redirect(302, "/gui/repository/imageinformation/list")
		guimessage.RedirectMessage(c)
		return
	}

	// Generate random work space
	workingDirectory := "/tmp/tmp_" + random.UUID()

	buildParameter := make(map[string]string)
	buildParameter["workingDirectory"] = workingDirectory
	buildParameter["repositoryPath"] = c.GetString("repositoryPath")
	buildParameter["sourceCodeProject"] = c.GetString("sourceCodeProject")
	buildParameter["sourceCodeDirectory"] = c.GetString("sourceCodeDirectory")
	buildParameter["sourceCodeMakeScript"] = c.GetString("sourceCodeMakeScript")
	buildParameter["environmentFile"] = c.GetString("environmentFile")

	switch kind {
	case "git":
		buildParameter["sourceCodeURL"] = c.GetString("sourceCodeURL")
		buildParameter["versionFile"] = c.GetString("versionFile")
	case "scp":
		buildParameter["hostAndPort"] = c.GetString("hostAndPort")
		buildParameter["username"] = c.GetString("username")
		buildParameter["password"] = c.GetString("password")
		buildParameter["sourcePath"] = c.GetString("sourcePath")
		buildParameter["compressFileName"] = c.GetString("compressFileName")
		buildParameter["unpackageCommand"] = c.GetString("unpackageCommand")
		buildParameter["versionFile"] = c.GetString("versionFile")
	case "sftp":
		buildParameter["hostAndPort"] = c.GetString("hostAndPort")
		buildParameter["username"] = c.GetString("username")
		buildParameter["password"] = c.GetString("password")
		buildParameter["sourcePath"] = c.GetString("sourcePath")
		buildParameter["versionFile"] = c.GetString("versionFile")
	}

	imageInformation := ImageInformation{
		name,
		kind,
		description,
		"",
		buildParameter,
		"",
		"",
		"",
		"",
		"",
		"",
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imageinformations/create/"

	resultJsonMap := make(map[string]interface{})

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestPostWithStructure(url, imageInformation, &resultJsonMap, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		guimessage.AddSuccess("The build " + name + " is launched asynchronizedly")
	}

	c.Ctx.Redirect(302, "/gui/repository/imageinformation/log?imageInformation="+imageInformation.Name)

	guimessage.RedirectMessage(c)
}
