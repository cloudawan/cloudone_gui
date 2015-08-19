package deploybluegreen

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"strconv"
)

type SelectController struct {
	beego.Controller
}

func (c *SelectController) Get() {
	c.TplNames = "deploy/deploybluegreen/select.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort, _ := beego.AppConfig.Int("kubeapiPort")

	imageInformation := c.GetString("imageInformation")
	currentEnvironment := c.GetString("currentEnvironment")
	nodePort := c.GetString("nodePort")
	description := c.GetString("description")

	c.Data["actionButtonValue"] = "Update"
	c.Data["pageHeader"] = "Update Blue Green Deployment"
	c.Data["imageInformation"] = imageInformation

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/deploybluegreens/deployable/" + imageInformation + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	namespaceSlice := make([]string, 0)
	_, err := restclient.RequestGetWithStructure(url, &namespaceSlice)
	if err != nil {
		guimessage.AddDanger("Fail to get deployable namespace")
	} else {
		c.Data["namespaceSlice"] = namespaceSlice
		c.Data["currentEnvironment"] = currentEnvironment
		c.Data["nodePort"] = nodePort
		c.Data["description"] = description
	}

	guimessage.OutputMessage(c.Data)
}

func (c *SelectController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort, _ := beego.AppConfig.Int("kubeapiPort")

	imageInformation := c.GetString("imageInformation")
	namespace := c.GetString("namespace")
	nodePort, _ := c.GetInt("nodePort")
	description := c.GetString("description")
	sessionAffinity := c.GetString("sessionAffinity")

	deployBlueGreen := DeployBlueGreen{
		imageInformation,
		namespace,
		nodePort,
		description,
		sessionAffinity,
	}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/deploybluegreens/" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	_, err := restclient.RequestPutWithStructure(url, deployBlueGreen, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Create blue green deployment " + imageInformation + " success")
	}

	c.Ctx.Redirect(302, "/gui/deploy/deploybluegreen/")

	guimessage.RedirectMessage(c)
}
