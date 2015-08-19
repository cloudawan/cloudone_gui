package deploy

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"strconv"
	"strings"
)

type UpdateController struct {
	beego.Controller
}

type DeployUpdateInput struct {
	ImageInformationName string
	Version              string
	Description          string
	EnvironmentSlice     []ReplicationControllerContainerEnvironment
}

func (c *UpdateController) Get() {
	c.TplNames = "deploy/deploy/update.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	name := c.GetString("name")
	oldVersion := c.GetString("oldVersion")

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/imagerecords/" + name

	imageRecordSlice := make([]ImageRecord, 0)

	_, err := restclient.RequestGetWithStructure(url, &imageRecordSlice)
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		versionEnvironmentMap := make(map[string]map[string]string)
		versionSlice := make([]string, 0)
		for _, imageRecord := range imageRecordSlice {
			if imageRecord.Version != oldVersion {
				versionSlice = append(versionSlice, imageRecord.Version)
				versionEnvironmentMap[imageRecord.Version] = imageRecord.Environment
			}
		}

		c.Data["name"] = name
		c.Data["versionSlice"] = versionSlice
		c.Data["versionEnvironmentMap"] = versionEnvironmentMap
	}

	guimessage.OutputMessage(c.Data)
}

func (c *UpdateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort, _ := beego.AppConfig.Int("kubeapiPort")

	namespaces, _ := c.GetSession("namespace").(string)

	imageInformationName := c.GetString("name")
	version := c.GetString("version")
	description := c.GetString("description")

	keySlice := make([]string, 0)
	inputMap := c.Input()
	if inputMap != nil {
		for key, _ := range inputMap {
			// Only collect environment belonging to this version
			if strings.HasPrefix(key, version) {
				keySlice = append(keySlice, key)
			}
		}
	}

	environmentSlice := make([]ReplicationControllerContainerEnvironment, 0)
	length := len(version) + 1 // + 1 for _
	for _, key := range keySlice {
		value := c.GetString(key)
		if len(value) > 0 {
			environmentSlice = append(environmentSlice,
				ReplicationControllerContainerEnvironment{key[length:], value})
		}
	}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/deploys/update/" + namespaces + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	deployUpdateInput := DeployUpdateInput{imageInformationName, version, description, environmentSlice}

	_, err := restclient.RequestPutWithStructure(url, deployUpdateInput, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Update deploy " + imageInformationName + " to version " + version + " success")
	}

	// Redirect to list
	c.Ctx.Redirect(302, "/gui/deploy/deploy/")

	guimessage.RedirectMessage(c)
}
