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

package user

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"strings"
)

type EditController struct {
	beego.Controller
}

type Role struct {
	Name        string
	Description string
	TagChecked  string
}

type Namespace struct {
	Name       string
	TagChecked string
}

func (c *EditController) Get() {
	c.TplName = "system/rbac/user/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	action := c.GetString("action")
	name := c.GetString("name")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		c.Ctx.Redirect(302, "/gui/system/rbac/user/list")
		guimessage.RedirectMessage(c)
		return
	}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/authorizations/roles"

	roleSlice := make([]Role, 0)

	_, err = restclient.RequestGetWithStructure(url, &roleSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		c.Ctx.Redirect(302, "/gui/system/rbac/user/list")
		guimessage.RedirectMessage(c)
		return
	}

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/namespaces" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	namespaceNameSlice := make([]string, 0)

	_, err = restclient.RequestGetWithStructure(url, &namespaceNameSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		c.Ctx.Redirect(302, "/gui/system/rbac/user/list")
		guimessage.RedirectMessage(c)
		return
	}

	namespaceSlice := make([]Namespace, 0)
	namespaceSlice = append(namespaceSlice, Namespace{"*", ""})
	for _, namespaceName := range namespaceNameSlice {
		namespaceSlice = append(namespaceSlice, Namespace{namespaceName, ""})
	}

	c.Data["action"] = action
	c.Data["roleSlice"] = roleSlice
	c.Data["namespaceSlice"] = namespaceSlice

	if action == "create" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Service"
	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/authorizations/users/" + name

		user := rbac.User{}

		_, err = restclient.RequestGetWithStructure(url, &user, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
			c.Ctx.Redirect(302, "/gui/system/rbac/user/list")
			guimessage.RedirectMessage(c)
			return
		}

		for i := 0; i < len(roleSlice); i++ {
			for _, ownedRole := range user.RoleSlice {
				if roleSlice[i].Name == ownedRole.Name {
					roleSlice[i].TagChecked = "checked"
				}
			}
		}

		for _, ownedResource := range user.ResourceSlice {
			// Simplied user version so the component of the resource is *
			if ownedResource.Path == "*" {
				// The first one is * all
				namespaceSlice[0].TagChecked = "checked"
			} else if strings.HasPrefix(ownedResource.Path, "/namespaces/") {
				splitSlice := strings.Split(ownedResource.Path, "/")
				ownedNamespaceName := splitSlice[2]
				for i := 0; i < len(namespaceSlice); i++ {
					if namespaceSlice[i].Name == ownedNamespaceName {
						namespaceSlice[i].TagChecked = "checked"
					}
				}
			}
		}

		c.Data["name"] = name
		c.Data["description"] = user.Description

		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Service"
	}

	guimessage.OutputMessage(c.Data)
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	name := c.GetString("name")
	password := c.GetString("password")
	description := c.GetString("description")
	action := c.GetString("action")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	roleSlice := make([]*rbac.Role, 0)
	resourceSlice := make([]*rbac.Resource, 0)

	namespaceNameSlice := make([]string, 0)
	hasNamespaceNameAll := false

	inputMap := c.Input()
	if inputMap != nil {
		for key, value := range inputMap {
			if strings.HasPrefix(key, "role_") {
				roleName := key[len("role_"):]
				if value[0] == "on" {
					roleSlice = append(roleSlice, &rbac.Role{roleName, nil, ""})
				}
			}
			if strings.HasPrefix(key, "namespace_") {
				namespaceName := key[len("namespace_"):]
				if value[0] == "on" {
					namespaceNameSlice = append(namespaceNameSlice, namespaceName)
					if namespaceName == "*" {
						hasNamespaceNameAll = true
					}
				}
			}
		}
	}

	if hasNamespaceNameAll {
		resourceSlice = append(resourceSlice, &rbac.Resource{"namespace_*", "*", "/namespaces/*"})
	} else {
		for _, namespaceName := range namespaceNameSlice {
			resourceSlice = append(resourceSlice, &rbac.Resource{"namespace_" + namespaceName, "*", "/namespaces/" + namespaceName})
		}
	}

	user := rbac.User{
		name,
		password,
		roleSlice,
		resourceSlice,
		description,
	}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	var err error
	if action == "create" {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/authorizations/users"

		_, err = restclient.RequestPostWithStructure(url, user, nil, tokenHeaderMap)

	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/authorizations/users/" + name

		_, err = restclient.RequestPutWithStructure(url, user, nil, tokenHeaderMap)
	}

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("User " + name + " is created")
	}

	c.Ctx.Redirect(302, "/gui/system/rbac/user/list")

	guimessage.RedirectMessage(c)
}
