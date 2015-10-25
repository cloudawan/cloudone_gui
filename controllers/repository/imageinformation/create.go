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
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/random"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
)

type CreateController struct {
	beego.Controller
}

func (c *CreateController) Get() {
	c.TplNames = "repository/imageinformation/create.html"
}

func (c *CreateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	name := c.GetString("name")
	kind := c.GetString("kind")
	description := c.GetString("description")

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
	}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/imageinformations/create/"

	_, err := restclient.RequestPostWithStructure(url, imageInformation, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Create  " + name + " success")
	}

	c.Ctx.Redirect(302, "/gui/repository/imageinformation/")

	guimessage.RedirectMessage(c)
}
