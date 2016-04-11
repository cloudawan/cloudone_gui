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

package deploy

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"strconv"
	"strings"
)

type CreateController struct {
	beego.Controller
}

type DeployCreateInput struct {
	ImageInformationName string
	Version              string
	Description          string
	ReplicaAmount        int
	PortSlice            []DeployContainerPort
	EnvironmentSlice     []ReplicationControllerContainerEnvironment
	ResourceMap          map[string]interface{}
}

type DeployContainerPort struct {
	Name          string
	ContainerPort int
	NodePort      int
}

type ReplicationControllerContainerPort struct {
	Name          string
	ContainerPort int
}

type ReplicationControllerContainerEnvironment struct {
	Name  string
	Value string
}

type ImageRecord struct {
	ImageInformation string
	Version          string
	Path             string
	VersionInfo      map[string]string
	Environment      map[string]string
	Description      string
	CreatedTime      string
}

type ByImageRecord []ImageRecord

func (b ByImageRecord) Len() int           { return len(b) }
func (b ByImageRecord) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByImageRecord) Less(i, j int) bool { return b[i].Version > b[j].Version } // Use > to list from latest to oldest

func (c *CreateController) Get() {
	c.TplName = "deploy/deploy/create.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/imagerecords/" + name

	imageRecordSlice := make([]ImageRecord, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &imageRecordSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		sort.Sort(ByImageRecord(imageRecordSlice))

		c.Data["imageInformationName"] = name
		c.Data["imageRecordSlice"] = imageRecordSlice
	}

	guimessage.OutputMessage(c.Data)
}

func (c *CreateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/deploy/deploy/list")
		return
	}

	namespaces, _ := c.GetSession("namespace").(string)

	imageInformationName := c.GetString("imageInformationName")
	version := c.GetString("version")
	description := c.GetString("description")
	replicaAmount, _ := c.GetInt("replicaAmount")

	resourceCPURequest, resourceCPURequestError := c.GetFloat("resourceCPURequest")
	resourceCPULimit, resourceCPULimitError := c.GetFloat("resourceCPULimit")
	resourceMemoryRequest, resourceMemoryRequestError := c.GetInt("resourceMemoryRequest")
	resourceMemoryLimit, resourceMemoryLimitError := c.GetInt("resourceMemoryLimit")

	portName := "generated"

	indexContainerPortMap := make(map[string]int)

	keySlice := make([]string, 0)
	inputMap := c.Input()
	if inputMap != nil {
		for key, _ := range inputMap {
			// Only collect environment belonging to this version
			if strings.HasPrefix(key, version) {
				keySlice = append(keySlice, key)
			}

			// Collect dynamically generated containerPort
			if strings.HasPrefix(key, "containerPort") {
				containerPort, _ := c.GetInt(key)
				index := key[len("containerPort"):]
				indexContainerPortMap[index] = containerPort
			}
		}
	}

	environmentSlice := make([]ReplicationControllerContainerEnvironment, 0)
	length := len(version) + 1 // + 1 for _
	for _, key := range keySlice {
		value := c.GetString(key)
		if len(value) > 0 {
			environmentSlice = append(environmentSlice,
				ReplicationControllerContainerEnvironment{key[length:], value})
		}
	}

	deployContainerPortSlice := make([]DeployContainerPort, 0)
	i := 0
	for index, containerPort := range indexContainerPortMap {
		nodePort := -1
		if c.GetString("useNodePort"+index) == "on" {
			if c.GetString("autoGeneratedNodePort"+index) == "on" {
				nodePort = 0
			} else {
				// If fail to parse, nodePort will be set to 0 that means auto-generated
				nodePort, _ = c.GetInt("nodePort" + index)
			}
		}
		deployContainerPortSlice = append(deployContainerPortSlice, DeployContainerPort{portName + strconv.Itoa(i), containerPort, nodePort})
		i++
	}

	// Resource reservation
	resourceMap := make(map[string]interface{})
	if resourceCPURequestError == nil || resourceMemoryRequestError == nil {
		resourceMap["requests"] = make(map[string]interface{})
		if resourceCPURequestError == nil {
			resourceMap["requests"].(map[string]interface{})["cpu"] = resourceCPURequest
		}
		if resourceMemoryRequestError == nil {
			resourceMap["requests"].(map[string]interface{})["memory"] = strconv.Itoa(resourceMemoryRequest) + "Mi"
		}
	}

	if resourceCPULimitError == nil || resourceMemoryLimitError == nil {
		resourceMap["limits"] = make(map[string]interface{})
		if resourceCPULimitError == nil {
			resourceMap["limits"].(map[string]interface{})["cpu"] = resourceCPULimit
		}
		if resourceMemoryLimitError == nil {
			resourceMap["limits"].(map[string]interface{})["memory"] = strconv.Itoa(resourceMemoryLimit) + "Mi"
		}
	}

	deployCreateInput := DeployCreateInput{
		imageInformationName,
		version,
		description,
		replicaAmount,
		deployContainerPortSlice,
		environmentSlice,
		resourceMap,
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/create/" + namespaces + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPostWithStructure(url, deployCreateInput, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Create deploy " + imageInformationName + " version " + version + " success")
	}

	c.Ctx.Redirect(302, "/gui/deploy/deploy/list")

	guimessage.RedirectMessage(c)
}
