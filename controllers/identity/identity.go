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

package identity

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/audit"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
)

const (
	loginPageURL  = "/gui/login"
	logoutPageURL = "/gui/logout"
)

func FilterUser(ctx *context.Context) {
	if (ctx.Input.IsGet() || ctx.Input.IsPost()) && (ctx.Input.URL() == loginPageURL || ctx.Input.URL() == logoutPageURL) {
		// Don't redirect itself to prevent the circle
	} else {
		user, ok := ctx.Input.Session("user").(*rbac.User)

		if ok == false {
			if guiMessage := guimessagedisplay.GetGUIMessageFromContext(ctx); guiMessage != nil {
				guiMessage.AddDanger("Username or password is incorrect")
			}
			ctx.Redirect(302, loginPageURL)
		} else {
			// Authorize
			if user.HasPermission(componentName, ctx.Input.Method(), ctx.Input.URL()) == false {
				if guiMessage := guimessagedisplay.GetGUIMessageFromContext(ctx); guiMessage != nil {
					guiMessage.AddDanger("User is not authorized to this page. Please use another user with priviledge.")
				}
				ctx.Redirect(302, loginPageURL)
			}

			// Resource check is in another place since GUI doesn't place the resource name in url

			// Audit log
			go func() {
				sendAuditLog(ctx, user.Name, true)
			}()
		}
	}
}

func sendAuditLog(ctx *context.Context, userName string, saveParameter bool) {
	cloudoneAnalysisProtocol := beego.AppConfig.String("cloudoneAnalysisProtocol")
	cloudoneAnalysisHost := beego.AppConfig.String("cloudoneAnalysisHost")
	cloudoneAnalysisPort := beego.AppConfig.String("cloudoneAnalysisPort")

	tokenHeaderMap, tokenHeaderMapOK := ctx.Input.Session("tokenHeaderMap").(map[string]string)
	requestURI := ctx.Input.URI()
	method := ctx.Input.Method()
	path := ctx.Input.URL()
	remoteAddress := ctx.Request.RemoteAddr
	queryParameterMap := ctx.Request.Form

	if saveParameter == false {
		// Not to save parameter, such as password
		requestURI = path
		queryParameterMap = nil
	}

	// Header is not used since the header has no useful information for now
	// Body is not used since the backend component will record again.
	// Path is not used since the backend component will record again.
	auditLog := audit.CreateAuditLog(componentName, path, userName, remoteAddress, queryParameterMap, nil, method, requestURI, "", nil)

	if tokenHeaderMapOK {
		url := cloudoneAnalysisProtocol + "://" + cloudoneAnalysisHost + ":" + cloudoneAnalysisPort + "/api/v1/auditlogs"

		_, err := restclient.RequestPost(url, auditLog, tokenHeaderMap, false)
		if err != nil {
			if guiMessage := guimessagedisplay.GetGUIMessageFromContext(ctx); guiMessage != nil {
				guiMessage.AddDanger("Fail to send audit log with error " + err.Error())
			}
		}
	}
}
