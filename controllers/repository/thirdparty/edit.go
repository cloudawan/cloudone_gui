package thirdparty

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"github.com/ghodss/yaml"
	"regexp"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplName = "repository/thirdparty/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")
	if name == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create third party service"
		c.Data["name"] = ""

		guimessage.OutputMessage(c.Data)
	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/clusterapplications/" + name
		cluster := Cluster{}

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestGetWithStructure(url, &cluster, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			guimessage.AddDanger("Fail to get with error" + err.Error())
			// Redirect to list
			c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")

			guimessage.RedirectMessage(c)
		} else {
			environmentByteSlice, err := json.MarshalIndent(cluster.Environment, "", "    ")
			if err != nil {
				guimessage.AddDanger("Fail to get with error" + err.Error())
				// Redirect to list
				c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")

				guimessage.RedirectMessage(c)
			}

			c.Data["actionButtonValue"] = "Update"
			c.Data["pageHeader"] = "Update third party service"
			c.Data["name"] = cluster.Name
			c.Data["description"] = cluster.Description
			c.Data["replicationControllerJson"] = string(cluster.ReplicationControllerJson)
			c.Data["serviceJson"] = cluster.ServiceJson
			c.Data["environment"] = string(environmentByteSlice)
			c.Data["scriptContent"] = cluster.ScriptContent

			switch cluster.ScriptType {
			case "none":
				c.Data["scriptTypeNone"] = "selected"
			case "python":
				c.Data["scriptTypePython"] = "selected"
			default:
			}

			guimessage.OutputMessage(c.Data)
		}
	}
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	name := c.GetString("name")
	description := c.GetString("description")

	// Name need to be a DNS 952 label
	match, _ := regexp.MatchString("^[a-z]{1}[a-z0-9-]{1,23}$", name)
	if match == false {
		guimessage.AddDanger("The name need to be a DNS 952 label ^[a-z]{1}[a-z0-9-]{1,23}$")
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")
		guimessage.RedirectMessage(c)
		return
	}

	replicationControllerJson := c.GetString("replicationControllerJson")
	if replicationControllerJson == "" {
		replicationControllerJson = "{}"
	}
	serviceJson := c.GetString("serviceJson")
	if serviceJson == "" {
		serviceJson = "{}"
	}
	environmentText := c.GetString("environment")
	if environmentText == "" {
		environmentText = "{}"
	}
	scriptType := c.GetString("scriptType")
	scriptContent := c.GetString("scriptContent")

	// Test replication controller
	replicationControllerJsonMap := make(map[string]interface{})
	// Try json
	err := json.Unmarshal([]byte(replicationControllerJson), &replicationControllerJsonMap)
	if err != nil {
		// Try yaml
		err := yaml.Unmarshal([]byte(replicationControllerJson), &replicationControllerJsonMap)
		if err != nil {
			// Error
			guimessage.AddDanger("Replication controller can't be parsed by json or yaml")
			guimessage.AddDanger(err.Error())
			c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")
			guimessage.RedirectMessage(c)
			return
		}
	}

	// Test service
	serviceJsonMap := make(map[string]interface{})
	// Try json
	err = json.Unmarshal([]byte(serviceJson), &serviceJsonMap)
	if err != nil {
		// Try yaml
		err := yaml.Unmarshal([]byte(serviceJson), &serviceJsonMap)
		if err != nil {
			// Error
			guimessage.AddDanger("Service can't be parsed by json or yaml")
			guimessage.AddDanger(err.Error())
			c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")
			guimessage.RedirectMessage(c)
			return
		}
	}

	environmentJsonMap := make(map[string]string)
	// Try json
	err = json.Unmarshal([]byte(environmentText), &environmentJsonMap)
	if err != nil {
		// Try yaml
		err := yaml.Unmarshal([]byte(environmentText), &environmentJsonMap)
		if err != nil {
			// Error
			guimessage.AddDanger("Environment can't be parsed by json or yaml")
			guimessage.AddDanger(err.Error())
			c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")
			guimessage.RedirectMessage(c)
			return
		}
	}

	cluster := Cluster{
		name,
		description,
		replicationControllerJson,
		serviceJson,
		environmentJsonMap,
		scriptType,
		scriptContent,
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/clusterapplications/"

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPost(url, cluster, tokenHeaderMap, true)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Third party application " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/repository/thirdparty/list")

	guimessage.RedirectMessage(c)
}
