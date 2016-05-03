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

package clone

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Clone struct {
	Name                           string
	DeployClusterApplication       DeployClusterApplication
	DeployInformation              DeployInformation
	RegionSlice                    []Region
	Order                          int
	CreatedTime                    time.Time
	DeployClusterApplicationHidden string
	DeployInformationHidden        string
}

type ByClone []Clone

func (b ByClone) Len() int           { return len(b) }
func (b ByClone) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByClone) Less(i, j int) bool { return b[i].CreatedTime.Before(b[j].CreatedTime) }

type DeployClusterApplication struct {
	Name                              string
	Namespace                         string
	Size                              int
	EnvironmentSlice                  []interface{}
	ReplicationControllerExtraJsonMap map[string]interface{}
	ServiceName                       string
	ReplicationControllerNameSlice    []string
	CreatedTime                       time.Time
	EnvironmentMap                    map[string]string
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
	Name     string
	Address  string
	Capacity Capacity
}

type Capacity struct {
	Cpu    string
	Memory string
}

type DeployInformation struct {
	Namespace                 string
	ImageInformationName      string
	CurrentVersion            string
	CurrentVersionDescription string
	Description               string
	ReplicaAmount             int
	ContainerPortSlice        []DeployContainerPort
	EnvironmentSlice          []ReplicationControllerContainerEnvironment
	ResourceMap               map[string]interface{}
	ExtraJsonMap              map[string]interface{}
	CreatedTime               time.Time
}

type DeployContainerPort struct {
	Name          string
	ContainerPort int
	NodePort      int
}

type ReplicationControllerContainerEnvironment struct {
	Name  string
	Value string
}

type Topology struct {
	Name            string
	SourceNamespace string
	CreatedUser     string
	CreatedDate     time.Time
	Description     string
	LaunchSlice     []Launch
}

type Launch struct {
	Order                    int
	LaunchApplication        *LaunchApplication
	LaunchClusterApplication *LaunchClusterApplication
}

type ByLaunch []Launch

func (b ByLaunch) Len() int           { return len(b) }
func (b ByLaunch) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByLaunch) Less(i, j int) bool { return b[i].Order < b[j].Order }

type LaunchApplication struct {
	ImageInformationName string
	Version              string
	Description          string
	ReplicaAmount        int
	PortSlice            []DeployContainerPort
	EnvironmentSlice     []ReplicationControllerContainerEnvironment
	ResourceMap          map[string]interface{}
	ExtraJsonMap         map[string]interface{}
}

type LaunchClusterApplication struct {
	Name                              string
	Size                              int
	EnvironmentSlice                  []interface{}
	ReplicationControllerExtraJsonMap map[string]interface{}
}

type TopologyController struct {
	beego.Controller
}

func (c *TopologyController) Get() {
	c.TplName = "deploy/clone/topology.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.OutputMessage(c.Data)
		return
	}

	currentNamespace, _ := c.GetSession("namespace").(string)

	action := c.GetString("action")
	c.Data["action"] = action
	if action == "clone" {
		c.Data["hiddenTagRegionTemplate"] = "hidden"
		c.Data["templateNameRequired"] = ""
		c.Data["actionButtonValue"] = "Clone"
	} else if action == "template" {
		c.Data["hiddenTagRegionTemplate"] = ""
		c.Data["templateNameRequired"] = "required"
		c.Data["actionButtonValue"] = "Create Template"
	}

	namespace := c.GetString("namespace")
	c.Data["sourceNamespace"] = namespace

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/nodes/topology?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	regionSlice := make([]Region, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestGetWithStructure(url, &regionSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger("Fail to get node topology with error" + err.Error())
		guimessage.OutputMessage(c.Data)
		return
	}

	filteredRegionSlice := make([]Region, 0)
	for _, region := range regionSlice {
		if region.LocationTagged {
			filteredRegionSlice = append(filteredRegionSlice, region)
		}
	}

	// Third-party service
	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deployclusterapplications/" + namespace

	deployClusterApplicationSlice := make([]DeployClusterApplication, 0)

	_, err = restclient.RequestGetWithStructure(url, &deployClusterApplicationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.OutputMessage(c.Data)
		return
	}

	for i, deployClusterApplication := range deployClusterApplicationSlice {
		if deployClusterApplication.EnvironmentSlice != nil {
			environmentMap := make(map[string]string)
			for _, environment := range deployClusterApplication.EnvironmentSlice {
				environmentJsonMap, ok := environment.(map[string]interface{})
				if ok {
					name, nameOk := environmentJsonMap["name"].(string)
					value, valueOk := environmentJsonMap["value"].(string)

					if nameOk && valueOk {
						// Try to set the known common parameter
						if name == "SERVICE_NAME" {
							value = deployClusterApplication.Name
						}
						if name == "NAMESPACE" {
							value = currentNamespace
						}
						if name == "GLUSTERFS_PATH_LIST" {
							value = ""
						}

						environmentMap[name] = value
					}
				}
			}

			deployClusterApplicationSlice[i].EnvironmentMap = environmentMap
		}
	}

	// Application
	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/" + namespace

	deployInformationSlice := make([]DeployInformation, 0)

	_, err = restclient.RequestGetWithStructure(url, &deployInformationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.OutputMessage(c.Data)
		return
	}

	for i, deployInformation := range deployInformationSlice {
		if deployInformation.EnvironmentSlice != nil {
			// Try to set the known common parameter
			for j, environment := range deployInformation.EnvironmentSlice {
				if environment.Name == "SERVICE_NAME" {
					deployInformationSlice[i].EnvironmentSlice[j].Value = deployInformation.ImageInformationName
				}
				if environment.Name == "NAMESPACE" {
					deployInformationSlice[i].EnvironmentSlice[j].Value = currentNamespace
				}
			}
		}
	}

	cloneSlice := make([]Clone, 0)
	for _, deployClusterApplication := range deployClusterApplicationSlice {
		clone := Clone{
			deployClusterApplication.Name,
			deployClusterApplication,
			DeployInformation{},
			filteredRegionSlice,
			0,
			deployClusterApplication.CreatedTime,
			"",
			"hidden",
		}
		// To avoid client side input validation
		clone.DeployInformation.ReplicaAmount = 1
		cloneSlice = append(cloneSlice, clone)
	}

	for _, deployInformation := range deployInformationSlice {
		clone := Clone{
			deployInformation.ImageInformationName,
			DeployClusterApplication{},
			deployInformation,
			filteredRegionSlice,
			0,
			deployInformation.CreatedTime,
			"hidden",
			"",
		}
		// To avoid client side input validation
		clone.DeployClusterApplication.Size = 1
		cloneSlice = append(cloneSlice, clone)
	}

	sort.Sort(ByClone(cloneSlice))
	c.Data["cloneSlice"] = cloneSlice

	for i := 0; i < len(cloneSlice); i++ {
		cloneSlice[i].Order = i + 1
	}

	guimessage.OutputMessage(c.Data)
}

func (c *TopologyController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/deploy/clone/select")
		return
	}

	action := c.GetString("action")

	// Generate topology order
	inputMap := c.Input()
	cloneNameSlice := make([]string, 0)
	environmentMap := make(map[string]string)
	if inputMap != nil {
		for key, _ := range inputMap {
			// Collect the clone name
			if strings.HasPrefix(key, "cloneUse") {
				cloneName := key[len("cloneUse"):]
				if c.GetString(key) == "on" && len(cloneName) > 0 {
					cloneNameSlice = append(cloneNameSlice, cloneName)
				}
				// Collect environment
			} else if strings.HasPrefix(key, "clusterEnvironment") {
				environmentMap[key] = c.GetString(key)
			} else if strings.HasPrefix(key, "applicationEnvironment") {
				environmentMap[key] = c.GetString(key)
			}
		}
	}

	sourceNamespace := c.GetString("sourceNamespace")

	// Application
	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/" + sourceNamespace

	deployInformationSlice := make([]DeployInformation, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestGetWithStructure(url, &deployInformationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/deploy/clone/select")
		return
	}

	launchSlice := make([]Launch, 0)
	for _, cloneName := range cloneNameSlice {
		cloneOrder, _ := c.GetInt("cloneOrder" + cloneName)

		clusterName := c.GetString("clusterName" + cloneName)
		clusterSize, _ := c.GetInt("clusterSize" + cloneName)

		applicationImageInformationName := c.GetString("applicationImageInformationName" + cloneName)
		applicationVersion := c.GetString("applicationVersion" + cloneName)
		applicationDescription := c.GetString("applicationDescription" + cloneName)
		applicationReplicaAmount, _ := c.GetInt("applicationReplicaAmount" + cloneName)

		clusterEnvironmentMap := make(map[string]string)
		applicationEnvironmentMap := make(map[string]string)
		for key, value := range environmentMap {
			if strings.HasPrefix(key, "clusterEnvironment"+cloneName) {
				environemntKey := key[len("clusterEnvironment"+cloneName):]
				clusterEnvironmentMap[environemntKey] = value
			} else if strings.HasPrefix(key, "applicationEnvironment"+cloneName) {
				environemntKey := key[len("applicationEnvironment"+cloneName):]
				applicationEnvironmentMap[environemntKey] = value
			}
		}

		// Location Affinity
		region := c.GetString("clusterRegion" + cloneName)
		zone := c.GetString("clusterZone" + cloneName)

		if region == "Any" {
			region = ""
		}
		if zone == "Any" {
			zone = ""
		}

		extraJsonMap := make(map[string]interface{})
		if len(region) > 0 {
			extraJsonMap["spec"] = make(map[string]interface{})
			extraJsonMap["spec"].(map[string]interface{})["template"] = make(map[string]interface{})
			extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"] = make(map[string]interface{})
			extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["nodeSelector"] = make(map[string]interface{})

			extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["nodeSelector"].(map[string]interface{})["region"] = region
			if len(zone) > 0 {
				extraJsonMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["nodeSelector"].(map[string]interface{})["zone"] = zone
			}
		} else {
			extraJsonMap = nil
		}

		// Cluster or application
		if len(clusterName) > 0 {
			environmentSlice := make([]interface{}, 0)
			for key, value := range clusterEnvironmentMap {
				environmentJsonMap := make(map[string]interface{})
				environmentJsonMap["name"] = key
				environmentJsonMap["value"] = value
				environmentSlice = append(environmentSlice, environmentJsonMap)
			}

			clusterLaunch := &LaunchClusterApplication{
				clusterName,
				clusterSize,
				environmentSlice,
				extraJsonMap,
			}

			launch := Launch{
				cloneOrder,
				nil,
				clusterLaunch,
			}

			launchSlice = append(launchSlice, launch)
		} else if len(applicationImageInformationName) > 0 {
			for _, deployInformation := range deployInformationSlice {
				if deployInformation.ImageInformationName == applicationImageInformationName {
					environmentSlice := make([]ReplicationControllerContainerEnvironment, 0)
					for key, value := range applicationEnvironmentMap {
						environmentSlice = append(environmentSlice, ReplicationControllerContainerEnvironment{key, value})
					}

					// Change the assigned Node Port to auto generated
					for i, containerPort := range deployInformation.ContainerPortSlice {
						if containerPort.NodePort > 0 {
							deployInformation.ContainerPortSlice[i].NodePort = 0
						}
					}

					launchApplication := &LaunchApplication{
						applicationImageInformationName,
						applicationVersion,
						applicationDescription,
						applicationReplicaAmount,
						deployInformation.ContainerPortSlice,
						environmentSlice,
						deployInformation.ResourceMap,
						extraJsonMap,
					}

					launch := Launch{
						cloneOrder,
						launchApplication,
						nil,
					}

					launchSlice = append(launchSlice, launch)

					break
				}
			}
		}
	}

	sort.Sort(ByLaunch(launchSlice))

	// Action: clone or create tempalte
	if action == "clone" {
		namespace, _ := c.GetSession("namespace").(string)

		for _, launch := range launchSlice {
			if launch.LaunchApplication != nil {
				url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
					"/api/v1/deploys/create/" + namespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

				_, err = restclient.RequestPostWithStructure(url, launch.LaunchApplication, nil, tokenHeaderMap)

				if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
					return
				}

				if err != nil {
					// Error
					guimessage.AddDanger(err.Error())
					guimessage.RedirectMessage(c)
					c.Ctx.Redirect(302, "/gui/deploy/clone/select")
					return
				}
			}
			if launch.LaunchClusterApplication != nil {
				url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
					"/api/v1/clusterapplications/launch/" + namespace + "/" + launch.LaunchClusterApplication.Name +
					"?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)
				jsonMap := make(map[string]interface{})

				_, err = restclient.RequestPostWithStructure(url, launch.LaunchClusterApplication, &jsonMap, tokenHeaderMap)

				if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
					return
				}

				if err != nil {
					// Error
					errorMessage, _ := jsonMap["Error"].(string)
					if strings.HasPrefix(errorMessage, "Replication controller already exists") {
						guimessage.AddDanger("Replication controller " + launch.LaunchClusterApplication.Name + " already exists")
					} else {
						guimessage.AddDanger(err.Error())
					}

					guimessage.RedirectMessage(c)
					c.Ctx.Redirect(302, "/gui/deploy/clone/select")
					return
				}
			}
		}

		guimessage.AddSuccess("Clone the namespace " + sourceNamespace + " to the namespace " + namespace)
	} else if action == "template" {
		templateName := c.GetString("templateName")
		templateDescription := c.GetString("templateDescription")

		createdUserName := ""
		user, ok := c.GetSession("user").(*rbac.User)
		if ok {
			createdUserName = user.Name
		}

		topology := Topology{
			templateName,
			sourceNamespace,
			createdUserName,
			time.Now(),
			templateDescription,
			launchSlice,
		}

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/topology/"

		_, err = restclient.RequestPostWithStructure(url, topology, nil, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
			guimessage.RedirectMessage(c)
			c.Ctx.Redirect(302, "/gui/deploy/clone/select")
			return
		}

		guimessage.AddSuccess("Create topology template from the namespace " + sourceNamespace + " as the template " + templateName)
	} else {
		guimessage.AddDanger("No such action: " + action)
	}

	c.Ctx.Redirect(302, "/gui/deploy/clone/select")

	guimessage.RedirectMessage(c)
}
