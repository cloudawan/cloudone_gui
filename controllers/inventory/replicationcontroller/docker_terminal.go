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

package replicationcontroller

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/ioutility"
	"github.com/cloudawan/cloudone_utility/restclient"
	"github.com/cloudawan/cloudone_utility/sshclient"
	"golang.org/x/net/websocket"
	"io"
	"strconv"
	"time"
)

type TerminalController struct {
	beego.Controller
}

func (c *TerminalController) Get() {
	c.TplName = "inventory/replicationcontroller/docker_terminal.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneGUIHost := c.Ctx.Input.Host()
	cloudoneGUIPort := c.Ctx.Input.Port()

	hostIP := c.GetString("hostIP")
	containerID := c.GetString("containerID")

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)
	token, _ := tokenHeaderMap["token"]

	c.Data["cloudoneGUIHost"] = cloudoneGUIHost
	c.Data["cloudoneGUIPort"] = cloudoneGUIPort

	c.Data["hostIP"] = hostIP
	c.Data["containerID"] = containerID

	c.Data["token"] = token

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

type WebSocketController struct {
	beego.Controller
}

func (c *WebSocketController) Get() {
	server := websocket.Server{Handler: ProxyServer}
	server.ServeHTTP(c.Ctx.ResponseWriter, c.Ctx.Request)
}

func ProxyServer(ws *websocket.Conn) {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	parameterMap := ws.Request().URL.Query()
	widthSlice := parameterMap["width"]
	heightSlice := parameterMap["height"]
	hostIPSlice := parameterMap["hostIP"]
	containerIDSlice := parameterMap["containerID"]
	tokenSlice := parameterMap["token"]

	if len(widthSlice) != 1 {
		errorMessage := "Parameter width is incorrect"
		ws.Write([]byte(errorMessage))
		ws.Close()
		return
	}
	width, err := strconv.Atoi(widthSlice[0])
	if err != nil {
		errorMessage := "Format of parameter width is incorrect"
		ws.Write([]byte(errorMessage))
		ws.Close()
		return
	}

	if len(heightSlice) != 1 {
		errorMessage := "Parameter height is incorrect"
		ws.Write([]byte(errorMessage))
		ws.Close()
		return
	}
	height, err := strconv.Atoi(heightSlice[0])
	if err != nil {
		errorMessage := "Format of parameter height is incorrect"
		ws.Write([]byte(errorMessage))
		ws.Close()
		return
	}

	if len(hostIPSlice) != 1 {
		errorMessage := "Parameter hostIP is incorrect"
		ws.Write([]byte(errorMessage))
		ws.Close()
		return
	}
	hostIP := hostIPSlice[0]

	if len(containerIDSlice) != 1 {
		errorMessage := "Parameter containerID is incorrect"
		ws.Write([]byte(errorMessage))
		ws.Close()
		return
	}
	containerID := containerIDSlice[0]
	// Remove docker protocol prefix docker://
	containerID = containerID[9:]

	if len(tokenSlice) != 1 {
		errorMessage := "Parameter token is incorrect"
		ws.Write([]byte(errorMessage))
		ws.Close()
		return
	}
	token := tokenSlice[0]
	headerMap := make(map[string]string)
	headerMap["token"] = token

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/hosts/credentials/" + hostIP

	credential := Credential{}

	_, err = restclient.RequestGetWithStructure(url, &credential, headerMap)

	if identity.IsTokenInvalid(err) {
		ws.Write([]byte(err.Error()))
		ws.Close()
		return
	}

	if err != nil {
		ws.Write([]byte(err.Error()))
		ws.Close()
		return
	}

	interactiveMap := make(map[string]string)
	interactiveMap["[sudo]"] = credential.SSH.Password + "\n"
	interactiveMap["exit"] = "exit\n"

	// Command way
	sshCommandProxy := sshclient.CreateSSHCommandProxy(
		2*time.Second,
		credential.IP,
		credential.SSH.Port,
		credential.SSH.User,
		credential.SSH.Password,
		height,
		width,
		interactiveMap)
	i, o, _, err := sshCommandProxy.Connect()

	if err != nil {
		ws.Write([]byte(err.Error()))
		ws.Close()
		return
	}

	i <- "sudo docker exec -it " + containerID + " bash\n"

	go func() {
		for {
			data, ok := <-o
			if ok == false {
				break
			} else {
				ws.Write([]byte(data))
			}
		}
		ws.Close()
	}()

	for {
		data, _, err := ioutility.ReadText(ws, 256)
		if err == io.EOF {
			break
		} else if err != nil {
			ws.Write([]byte(err.Error()))
			break
		}

		i <- data
	}
	ws.Close()
	sshCommandProxy.Disconnect()

	// Stream way
	/*
		sshStreamProxy := sshclient.CreateSSHStreamProxy(
			2*time.Second,
			credential.IP,
			credential.SSH.Port,
			credential.SSH.User,
			credential.SSH.Password,
			height,
			width)
		w, r, _, err := sshStreamProxy.Connect()

		if err != nil {
			ws.Write([]byte(err.Error()))
			ws.Close()
			return
		}

		w.Write([]byte("sudo docker exec -it " + containerID + " bash\n"))
		w.Write([]byte(credential.SSH.Password + "\n"))

		go func() {
			// Cancatenate ssh reader to websocket writer
			io.Copy(ws, r)

			ws.Close()
		}()

		// Cancatenate websocket reader to ssh writer
		io.Copy(w, ws)

		ws.Close()
		sshStreamProxy.Disconnect()
	*/
}
