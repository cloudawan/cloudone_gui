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
	"sort"
	"strconv"
	"strings"
	"time"
)

type EditController struct {
	beego.Controller
}

type Role struct {
	Name        string
	Description string
	Tag         string
}

type Namespace struct {
	Name string
	Tag  string
}

type ByRole []Role

func (b ByRole) Len() int           { return len(b) }
func (b ByRole) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByRole) Less(i, j int) bool { return b[i].Name < b[j].Name }

type ByNamespace []Namespace

func (b ByNamespace) Len() int           { return len(b) }
func (b ByNamespace) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByNamespace) Less(i, j int) bool { return b[i].Name < b[j].Name }

const (
	guiWidgetTimePickerFormat = "01/02/2006 15:04 PM"
)

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

	loginNamespaceSlice := make([]Namespace, 0)
	namespaceSlice := make([]Namespace, 0)
	namespaceSlice = append(namespaceSlice, Namespace{"*", ""})
	for _, namespaceName := range namespaceNameSlice {
		namespaceSlice = append(namespaceSlice, Namespace{namespaceName, ""})
		loginNamespaceSlice = append(loginNamespaceSlice, Namespace{namespaceName, ""})
	}

	sort.Sort(ByRole(roleSlice))
	sort.Sort(ByNamespace(namespaceSlice))

	c.Data["action"] = action
	c.Data["roleSlice"] = roleSlice
	c.Data["namespaceSlice"] = namespaceSlice
	c.Data["loginNamespaceSlice"] = loginNamespaceSlice

	if action == "create" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create User"
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
					roleSlice[i].Tag = "checked"
				}
			}
		}

		for _, ownedResource := range user.ResourceSlice {
			// Simplied user version so the component of the resource is *
			if ownedResource.Path == "*" || ownedResource.Path == "/namespaces/" {
				// The first one is * all
				namespaceSlice[0].Tag = "checked"
			} else if strings.HasPrefix(ownedResource.Path, "/namespaces/") {
				splitSlice := strings.Split(ownedResource.Path, "/")
				ownedNamespaceName := splitSlice[2]

				for i := 0; i < len(namespaceSlice); i++ {
					if namespaceSlice[i].Name == ownedNamespaceName {
						namespaceSlice[i].Tag = "checked"
					}
				}
			}
		}

		metaDataMap := user.MetaDataMap
		if metaDataMap == nil {
			metaDataMap = make(map[string]string)
		}

		loginNamespace := metaDataMap["loginNamespace"]
		if len(loginNamespace) > 0 {
			for i := 0; i < len(loginNamespaceSlice); i++ {
				if loginNamespaceSlice[i].Name == loginNamespace {
					loginNamespaceSlice[i].Tag = "selected"
				}
			}
		}

		if user.Disabled {
			c.Data["disabledChecked"] = "checked"
		}

		if user.ExpiredTime != nil {
			c.Data["expiredTime"] = user.ExpiredTime.Format(guiWidgetTimePickerFormat)
		}

		c.Data["name"] = name
		c.Data["description"] = user.Description
		c.Data["githubWebhookSecret"] = user.MetaDataMap["githubWebhookSecret"]
		c.Data["readonly"] = "readonly"

		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update User"
	}

	guimessage.OutputMessage(c.Data)
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	name := c.GetString("name")
	password := c.GetString("password")
	disabledText := c.GetString("disabled")
	expiredTimeText := c.GetString("expiredTime")
	description := c.GetString("description")
	githubWebhookSecret := c.GetString("githubWebhookSecret")
	action := c.GetString("action")

	loginNamespace := c.GetString("loginNamespace")

	disabled := false
	if disabledText == "on" {
		disabled = true
	}

	var expiredTime *time.Time = nil
	expiredTimeData, err := time.Parse(guiWidgetTimePickerFormat, expiredTimeText)
	if err == nil {
		expiredTime = &expiredTimeData
	}

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
		resourceSlice = append(resourceSlice, &rbac.Resource{"namespace_*", "*", "/namespaces/"})
	} else {
		for _, namespaceName := range namespaceNameSlice {
			resourceSlice = append(resourceSlice, &rbac.Resource{"namespace_" + namespaceName, "*", "/namespaces/" + namespaceName})
		}
	}

	metaDataMap := make(map[string]string)
	if len(loginNamespace) > 0 {
		metaDataMap["loginNamespace"] = loginNamespace
	}

	if len(githubWebhookSecret) > 0 {
		metaDataMap["githubWebhookSecret"] = githubWebhookSecret
	}

	user := rbac.User{
		name,
		password,
		roleSlice,
		resourceSlice,
		description,
		metaDataMap,
		expiredTime,
		disabled,
	}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

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
		guimessage.AddSuccess("User " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/system/rbac/user/list")

	guimessage.RedirectMessage(c)
}
