package thirdparty

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"strings"
)

type Cluster struct {
	Name                      string
	Description               string
	ReplicationControllerJson string
	ServiceJson               string
	Environment               map[string]string
	ScriptType                string
	ScriptContent             string
}

type LaunchController struct {
	beego.Controller
}

func (c *LaunchController) Get() {
	c.TplNames = "repository/thirdparty/launch.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	name := c.GetString("name")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/clusterapplications/" + name
	cluster := Cluster{}
	_, err := restclient.RequestGetWithStructure(url, &cluster)

	if err != nil {
		guimessage.AddDanger("Fail to get with error" + err.Error())
		// Redirect to list
		c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

		guimessage.RedirectMessage(c)
	} else {
		c.Data["actionButtonValue"] = "Launch"
		c.Data["pageHeader"] = "Launch third party service"
		c.Data["thirdPartyApplicationName"] = name
		c.Data["environment"] = cluster.Environment

		guimessage.OutputMessage(c.Data)
	}
}

func (c *LaunchController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort := beego.AppConfig.String("kubeapiPort")

	namespace, _ := c.GetSession("namespace").(string)
	name := c.GetString("name")
	size := c.GetString("size")

	keySlice := make([]string, 0)
	inputMap := c.Input()
	if inputMap != nil {
		for key, _ := range inputMap {
			keySlice = append(keySlice, key)
		}
	}

	environmentSlice := make([]interface{}, 0)
	for _, key := range keySlice {
		value := c.GetString(key)
		if len(value) > 0 {
			environmentMap := make(map[string]string)
			environmentMap["name"] = key
			environmentMap["value"] = value
			environmentSlice = append(environmentSlice, environmentMap)
		}
	}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/clusterapplications/launch/" + namespace + "/" + name +
		"?kubeapihost=" + kubeapiHost + "&kubeapiport=" + kubeapiPort + "&size=" + size
	jsonMap := make(map[string]interface{})
	nil, err := restclient.RequestPostWithStructure(url, environmentSlice, &jsonMap)

	if err != nil {
		// Error
		errorMessage, _ := jsonMap["Error"].(string)
		if strings.HasPrefix(errorMessage, "Replication controller already exists") {
			guimessage.AddDanger("Replication controller " + name + " already exists")
		} else {
			guimessage.AddDanger(err.Error())
		}
	} else {
		guimessage.AddSuccess("Cluster application " + name + " is launched")
	}

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/repository/thirdparty/")

	guimessage.RedirectMessage(c)
}
