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

package upgrade

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"github.com/cloudawan/cloudone_utility/sshclient"
	"golang.org/x/net/websocket"
	"strconv"
	"time"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	c.TplName = "system/upgrade/index.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneGUIHost := c.Ctx.Input.Host()
	cloudoneGUIPort := c.Ctx.Input.Port()

	c.Data["cloudoneGUIHost"] = cloudoneGUIHost
	c.Data["cloudoneGUIPort"] = cloudoneGUIPort

	guimessage.OutputMessage(c.Data)
}

type Credential struct {
	IP  string
	SSH SSH
}

type SSH struct {
	Port     int
	User     string
	Password string
}

type ReplicationControllerAndRelatedPod struct {
	Name               string
	Namespace          string
	ReplicaAmount      int
	AliveReplicaAmount int
	Selector           map[string]string
	Label              map[string]string
	PodSlice           []Pod
	Display            string
}

type Pod struct {
	Name           string
	Namespace      string
	HostIP         string
	PodIP          string
	ContainerSlice []PodContainer
}

type PodContainer struct {
	Name        string
	Image       string
	ContainerID string
	PortSlice   []PodContainerPort
}

type PodContainerPort struct {
	Name          string
	ContainerPort int
	Protocol      string
}

type WebSocketController struct {
	beego.Controller
}

func (c *WebSocketController) Get() {
	server := websocket.Server{Handler: ProxyServer}
	server.ServeHTTP(c.Ctx.ResponseWriter, c.Ctx.Request)
}

func getParameter(parameterMap map[string][]string, name string) string {
	slice := parameterMap[name]
	if len(slice) == 1 {
		return slice[0]
	} else {
		return ""
	}
}

func UpgradeDockerImage(ws *websocket.Conn, credentialSlice []Credential, path string, version string) error {
	if len(credentialSlice) == 0 {
		return errors.New("No credential data")
	}

	if path == "" {
		return errors.New("Image path can't be empty")
	}

	imageUri := path
	if version != "" {
		imageUri = imageUri + ":" + version
	}

	totalAmount := len(credentialSlice)

	errorChannel := make(chan error, totalAmount)

	for _, credential := range credentialSlice {
		go func(credential Credential) {
			ws.Write([]byte("Start to pull " + imageUri + " on host " + credential.IP + "\n"))

			commandSlice := make([]string, 0)
			commandSlice = append(commandSlice, "sudo docker pull "+imageUri+"\n")
			interactiveMap := make(map[string]string)
			interactiveMap["[sudo]"] = credential.SSH.Password + "\n"

			resultSlice, err := sshclient.InteractiveSSH(
				2*time.Second,
				10*time.Minute,
				credential.IP,
				credential.SSH.Port,
				credential.SSH.User,
				credential.SSH.Password,
				commandSlice,
				interactiveMap)

			if err != nil {
				errorChannel <- err
				ws.Write([]byte(imageUri + " on host " + credential.IP + " has error to upgraded\n"))
				for _, result := range resultSlice {
					ws.Write([]byte(result + "\n"))
				}

			} else {
				errorChannel <- nil
				ws.Write([]byte(imageUri + " on host " + credential.IP + " is upgraded\n"))
			}
		}(credential)
	}

	// Wait for all go routine to finish
	hasError := false
	errorSlice := make([]error, 0)
	for i := 0; i < totalAmount; i++ {
		err, _ := <-errorChannel
		if err != nil {
			hasError = true
		}
		errorSlice = append(errorSlice, err)
	}

	close(errorChannel)

	if hasError {
		errorMessage := fmt.Sprintf("%v", errorSlice)
		return errors.New(errorMessage)
	} else {
		return nil
	}
}

func stopDockerContainer(credential Credential, containerID string) error {
	commandSlice := make([]string, 0)
	commandSlice = append(commandSlice, "sudo docker stop "+containerID+"\n")
	interactiveMap := make(map[string]string)
	interactiveMap["[sudo]"] = credential.SSH.Password + "\n"

	resultSlice, err := sshclient.InteractiveSSH(
		2*time.Second,
		3*time.Minute,
		credential.IP,
		credential.SSH.Port,
		credential.SSH.User,
		credential.SSH.Password,
		commandSlice,
		interactiveMap)

	if err != nil {
		errorMessage := fmt.Sprintf("%v %v", err, resultSlice)
		return errors.New(errorMessage)
	} else {
		return nil
	}
}

func ProxyServer(ws *websocket.Conn) {

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	cloudoneAnalysisProtocol := beego.AppConfig.String("cloudoneAnalysisProtocol")
	cloudoneAnalysisHost := beego.AppConfig.String("cloudoneAnalysisHost")
	cloudoneAnalysisPort := beego.AppConfig.String("cloudoneAnalysisPort")

	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()

	if err != nil {
		ws.Write([]byte(err.Error()))
		ws.Close()
		return
	}

	parameterMap := ws.Request().URL.Query()
	upgradeCloudone := getParameter(parameterMap, "upgradeCloudone")
	upgradeCloudoneImagePath := getParameter(parameterMap, "upgradeCloudoneImagePath")
	upgradeCloudoneVersion := getParameter(parameterMap, "upgradeCloudoneVersion")
	upgradeCloudoneGUI := getParameter(parameterMap, "upgradeCloudoneGUI")
	upgradeCloudoneGUIImagePath := getParameter(parameterMap, "upgradeCloudoneGUIImagePath")
	upgradeCloudoneGUIVersion := getParameter(parameterMap, "upgradeCloudoneGUIVersion")
	upgradeCloudoneAnalysis := getParameter(parameterMap, "upgradeCloudoneAnalysis")
	upgradeCloudoneAnalysisImagePath := getParameter(parameterMap, "upgradeCloudoneAnalysisImagePath")
	upgradeCloudoneAnalysisVersion := getParameter(parameterMap, "upgradeCloudoneAnalysisVersion")
	upgradeTopologyConfiguration := getParameter(parameterMap, "upgradeTopologyConfiguration")
	upgradeNamespace := getParameter(parameterMap, "upgradeNamespace")
	upgradeReplicationControllerName := getParameter(parameterMap, "upgradeReplicationControllerName")
	upgradeReplicationControllerContent := getParameter(parameterMap, "upgradeReplicationControllerContent")
	upgradeServiceName := getParameter(parameterMap, "upgradeServiceName")
	upgradeServiceContent := getParameter(parameterMap, "upgradeServiceContent")

	replicationControllerJsonMap := make(map[string]interface{})
	serviceJsonMap := make(map[string]interface{})
	if upgradeTopologyConfiguration == "true" {
		if upgradeReplicationControllerContent != "" {
			err := json.Unmarshal([]byte(upgradeReplicationControllerContent), &replicationControllerJsonMap)
			if err != nil {
				errorMessage := "Can't unmarshal upgradeReplicationControllerContent with error " + err.Error() + "\n"
				ws.Write([]byte(errorMessage))
				ws.Close()
				return
			}
		}
		if upgradeServiceContent != "" {
			err := json.Unmarshal([]byte(upgradeServiceContent), &serviceJsonMap)
			if err != nil {
				errorMessage := "Can't unmarshal upgradeServiceContent with error " + err.Error() + "\n"
				ws.Write([]byte(errorMessage))
				ws.Close()
				return
			}
		}
	}

	// Get credential
	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/hosts/credentials/"

	credentialSlice := make([]Credential, 0)

	_, err = restclient.RequestGetWithStructure(url, &credentialSlice)
	if err != nil {
		errorMessage := "Can't get credential data with error " + err.Error() + "\n"
		ws.Write([]byte(errorMessage))
		ws.Close()
		return
	}

	// Pull images
	if upgradeCloudone == "true" {
		err := UpgradeDockerImage(ws, credentialSlice, upgradeCloudoneImagePath, upgradeCloudoneVersion)
		if err != nil {
			errorMessage := "Can't upgrade image cloudone with error " + err.Error() + "\n"
			ws.Write([]byte(errorMessage))
			ws.Close()
			return
		}
	}

	if upgradeCloudoneGUI == "true" {
		err := UpgradeDockerImage(ws, credentialSlice, upgradeCloudoneGUIImagePath, upgradeCloudoneGUIVersion)
		if err != nil {
			errorMessage := "Can't upgrade image cloudone gui with error " + err.Error() + "\n"
			ws.Write([]byte(errorMessage))
			ws.Close()
			return
		}
	}

	if upgradeCloudoneAnalysis == "true" {
		err := UpgradeDockerImage(ws, credentialSlice, upgradeCloudoneAnalysisImagePath, upgradeCloudoneAnalysisVersion)
		if err != nil {
			errorMessage := "Can't upgrade image cloudone analysis with error " + err.Error() + "\n"
			ws.Write([]byte(errorMessage))
			ws.Close()
			return
		}
	}

	// Update service
	if len(serviceJsonMap) > 0 {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/services/json/" + upgradeNamespace + "/" + upgradeReplicationControllerName + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

		ws.Write([]byte("Start to update service\n"))

		_, err := restclient.RequestPutWithStructure(url, serviceJsonMap, nil)
		if err != nil {
			errorMessage := "Can't upgrade service with error " + err.Error() + "\n"
			ws.Write([]byte(errorMessage))
			ws.Close()
			return
		} else {
			ws.Write([]byte("The service is updated\n"))
		}
	}

	if len(replicationControllerJsonMap) > 0 {
		// Update replication controller to change all instances
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/replicationcontrollers/json/" + upgradeNamespace + "/" + upgradeServiceName + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

		ws.Write([]byte("Stop and recreate the replication controller. Please refresh the page after tens of seconds\n"))

		_, err := restclient.RequestPutWithStructure(url, replicationControllerJsonMap, nil)
		if err != nil {
			errorMessage := "Can't upgrade replication controller with error " + err.Error() + "\n"
			ws.Write([]byte(errorMessage))
			ws.Close()
			return
		} else {
			ws.Write([]byte("The replication controller is updated\n"))
		}
	} else {
		// Only update the specific docker containers
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/replicationcontrollers/" + upgradeNamespace + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

		replicationControllerAndRelatedPodSlice := make([]ReplicationControllerAndRelatedPod, 0)
		_, err := restclient.RequestGetWithStructure(url, &replicationControllerAndRelatedPodSlice)
		if err != nil {
			errorMessage := "Can't get replication controller and related pod data with error " + err.Error() + "\n"
			ws.Write([]byte(errorMessage))
			ws.Close()
			return
		}
		selectedReplicationControllerAndRelatedPod := ReplicationControllerAndRelatedPod{}
		for _, replicationControllerAndRelatedPod := range replicationControllerAndRelatedPodSlice {
			if replicationControllerAndRelatedPod.Name == upgradeReplicationControllerName {
				selectedReplicationControllerAndRelatedPod = replicationControllerAndRelatedPod
			}
		}

		for _, pod := range selectedReplicationControllerAndRelatedPod.PodSlice {
			usedCredential := Credential{}
			for _, credential := range credentialSlice {
				if credential.IP == pod.HostIP {
					usedCredential = credential
					break
				}
			}

			cloudoneContainer := PodContainer{}
			cloudoneAnalysisContainer := PodContainer{}
			cloudoneGUIContainer := PodContainer{}

			for _, container := range pod.ContainerSlice {
				if container.Name == "cloudone" {
					cloudoneContainer = container
				}
				if container.Name == "cloudone-analysis" {
					cloudoneAnalysisContainer = container
				}
				if container.Name == "cloudone-gui" {
					cloudoneGUIContainer = container
				}
			}

			if upgradeCloudone == "true" && cloudoneContainer.ContainerID != "" {
				containerID := cloudoneContainer.ContainerID[9:]
				ws.Write([]byte("Stop and recreate cloudone\n"))
				err := stopDockerContainer(usedCredential, containerID)
				if err != nil {
					errorMessage := "Can't stop container " + cloudoneContainer.Name + " with error " + err.Error() + "\n"
					ws.Write([]byte(errorMessage))
					ws.Close()
					return
				}
			}
			if upgradeCloudoneAnalysis == "true" && cloudoneAnalysisContainer.ContainerID != "" {
				containerID := cloudoneAnalysisContainer.ContainerID[9:]
				ws.Write([]byte("Stop and recreate cloudone_analysis\n"))
				err := stopDockerContainer(usedCredential, containerID)
				if err != nil {
					errorMessage := "Can't stop container " + cloudoneAnalysisContainer.Name + " with error " + err.Error() + "\n"
					ws.Write([]byte(errorMessage))
					ws.Close()
					return
				}
			}
			if upgradeCloudoneGUI == "true" && cloudoneGUIContainer.ContainerID != "" {
				containerID := cloudoneGUIContainer.ContainerID[9:]
				ws.Write([]byte("Stop and recreate cloudone_gui. Please refresh the page after tens of seconds\n"))
				err := stopDockerContainer(usedCredential, containerID)
				if err != nil {
					errorMessage := "Can't stop container " + cloudoneGUIContainer.Name + " with error " + err.Error() + "\n"
					ws.Write([]byte(errorMessage))
					ws.Close()
					return
				}
			}
		}

		if upgradeCloudone == "true" {
			url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
				"/api/v1/healthchecks/"
			cloudoneJsonMap := make(map[string]interface{}, 0)
			for {
				time.Sleep(time.Second)
				_, err := restclient.RequestGetWithStructure(url, &cloudoneJsonMap)
				if err == nil {
					break
				} else {
					ws.Write([]byte("Wait for Cloudone to come up\n"))
				}
			}
		}

		if upgradeCloudoneAnalysis == "true" {
			url := cloudoneAnalysisProtocol + "://" + cloudoneAnalysisHost + ":" + cloudoneAnalysisPort +
				"/api/v1/healthchecks/"
			cloudoneAnalysisJsonMap := make(map[string]interface{}, 0)
			for {
				time.Sleep(time.Second)
				_, err := restclient.RequestGetWithStructure(url, &cloudoneAnalysisJsonMap)
				if err == nil {
					break
				} else {
					ws.Write([]byte("Wait for Cloudone Analysis to come up\n"))
				}
			}
		}
	}

	ws.Write([]byte("Upgrade is done\n"))

	ws.Close()
}
