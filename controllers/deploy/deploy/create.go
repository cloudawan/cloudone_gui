package deploy

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"strconv"
	"strings"
)

type CreateController struct {
	beego.Controller
}

type DeployCreateInput struct {
	ImageInformationName string
	Version              string
	Description          string
	ReplicaAmount        int
	PortSlice            []ReplicationControllerContainerPort
	EnvironmentSlice     []ReplicationControllerContainerEnvironment
}

type ReplicationControllerContainerPort struct {
	Name          string
	ContainerPort int
}

type ReplicationControllerContainerEnvironment struct {
	Name  string
	Value string
}

type ImageRecord struct {
	ImageInformation string
	Version          string
	Path             string
	VersionInfo      map[string]string
	Environment      map[string]string
	Description      string
	CreatedTime      string
}

func (c *CreateController) Get() {
	c.TplNames = "deploy/deploy/create.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

	name := c.GetString("name")

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
			versionSlice = append(versionSlice, imageRecord.Version)
			versionEnvironmentMap[imageRecord.Version] = imageRecord.Environment
		}

		c.Data["imageInformationName"] = name
		c.Data["versionSlice"] = versionSlice
		c.Data["versionEnvironmentMap"] = versionEnvironmentMap
	}

	guimessage.OutputMessage(c.Data)
}

func (c *CreateController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort, _ := beego.AppConfig.Int("kubeapiPort")

	namespaces, _ := c.GetSession("namespace").(string)

	imageInformationName := c.GetString("imageInformationName")
	version := c.GetString("version")
	description := c.GetString("description")
	replicaAmount, _ := c.GetInt("replicaAmount")
	//portName := c.GetString("portName")
	containerPort, err := c.GetInt("containerPort")

	portName := "generated"

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

	replicationControllerContainerPortSlice := make([]ReplicationControllerContainerPort, 0)
	replicationControllerContainerPortSlice = append(replicationControllerContainerPortSlice, ReplicationControllerContainerPort{portName, containerPort})

	deployCreateInput := DeployCreateInput{
		imageInformationName,
		version,
		description,
		replicaAmount,
		replicationControllerContainerPortSlice,
		environmentSlice,
	}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/deploys/create/" + namespaces + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	_, err = restclient.RequestPostWithStructure(url, deployCreateInput, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Create deploy " + imageInformationName + " version " + version + " success")
	}

	c.Ctx.Redirect(302, "/gui/deploy/deploy/")

	guimessage.RedirectMessage(c)
}