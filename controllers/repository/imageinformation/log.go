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
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/hpcloud/tail"
	"golang.org/x/net/websocket"
	"os"
)

type LogController struct {
	beego.Controller
}

const (
	processOutMessageFilePathAndNamePrefix = "/tmp/processingBuildLog"
)

func GetProcessingOutMessageFilePathAndName(imageInformationName string) string {
	return processOutMessageFilePathAndNamePrefix + imageInformationName
}

func (c *LogController) Get() {
	c.TplName = "repository/imageinformation/log.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneGUIHost := c.Ctx.Input.Host()
	cloudoneGUIPort := c.Ctx.Input.Port()

	imageInformation := c.GetString("imageInformation")

	if _, err := os.Stat(GetProcessingOutMessageFilePathAndName(imageInformation)); os.IsNotExist(err) {
		// does not exist
		guimessage.AddInfo("No ongoing building process.")
		c.Ctx.Redirect(302, "/gui/repository/imageinformation/list")
		guimessage.RedirectMessage(c)
		return
	}

	c.Data["cloudoneGUIHost"] = cloudoneGUIHost
	c.Data["cloudoneGUIPort"] = cloudoneGUIPort
	c.Data["imageInformation"] = imageInformation

	guimessage.OutputMessage(c.Data)
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

func ProxyServer(ws *websocket.Conn) {
	parameterMap := ws.Request().URL.Query()
	imageInformation := getParameter(parameterMap, "imageInformation")

	filename := GetProcessingOutMessageFilePathAndName(imageInformation)

	// TODO The data should be from cloudon. Since they share the same file system now, read from here now.
	t, err := tail.TailFile(filename, tail.Config{Follow: true, Poll: true})
	if err != nil {
		ws.Write([]byte(err.Error()))
		ws.Close()
		return
	}

	for line := range t.Lines {
		ws.Write([]byte(line.Text))
		ws.Write([]byte("\n"))
	}

	t.Cleanup()

	ws.Close()
}
