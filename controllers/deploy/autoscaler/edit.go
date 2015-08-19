package autoscaler

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"time"
)

type EditController struct {
	beego.Controller
}

func (c *EditController) Get() {
	c.TplNames = "deploy/autoscaler/edit.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	kind := c.GetString("kind")
	name := c.GetString("name")

	c.Data["cpuHidden"] = "hidden"
	c.Data["memoryHidden"] = "hidden"

	if kind == "" || name == "" {
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Autoscaler"
		c.Data["kind"] = ""
		c.Data["name"] = ""
		c.Data["readonly"] = ""
	} else {
		namespace, _ := c.GetSession("namespace").(string)

		kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
		kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
		kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")

		url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
			"/api/v1/autoscalers/" + namespace + "/" + kind + "/" + name

		replicationControllerAutoScaler := ReplicationControllerAutoScaler{}

		_, err := restclient.RequestGetWithStructure(url, &replicationControllerAutoScaler)

		if err != nil {
			// Error
			guimessage.AddDanger(err.Error())
		} else {
			c.Data["maximumReplica"] = replicationControllerAutoScaler.MaximumReplica
			c.Data["minimumReplica"] = replicationControllerAutoScaler.MinimumReplica

			for _, indicator := range replicationControllerAutoScaler.IndicatorSlice {
				switch indicator.Type {
				case "cpu":
					c.Data["cpuChecked"] = "checked"
					delete(c.Data, "cpuHidden")

					if indicator.AboveAllOrOne {
						c.Data["cpuAboveAllOrOneChecked"] = "checked"
					}

					c.Data["cpuAbovePercentageOfData"] = int(indicator.AbovePercentageOfData * 100)
					c.Data["cpuAboveThreshold"] = indicator.AboveThreshold / 1000000

					if indicator.BelowAllOrOne {
						c.Data["cpuBelowAllOrOneChecked"] = "checked"
					}

					c.Data["cpuBelowPercentageOfData"] = int(indicator.BelowPercentageOfData * 100)
					c.Data["cpuBelowThreshold"] = indicator.BelowThreshold / 1000000
				case "memory":
					c.Data["memoryChecked"] = "checked"
					delete(c.Data, "memoryHidden")

					if indicator.AboveAllOrOne {
						c.Data["memoryAboveAllOrOneChecked"] = "checked"
					}

					c.Data["memoryAbovePercentageOfData"] = int(indicator.AbovePercentageOfData * 100)
					c.Data["memoryAboveThreshold"] = indicator.AboveThreshold / (1024 * 1024)

					if indicator.BelowAllOrOne {
						c.Data["memoryBelowAllOrOneChecked"] = "checked"
					}

					c.Data["memoryBelowPercentageOfData"] = int(indicator.BelowPercentageOfData * 100)
					c.Data["memoryBelowThreshold"] = indicator.BelowThreshold / (1024 * 1024)
				}
			}
			coolDownDurationInSecond := int(replicationControllerAutoScaler.CoolDownDuration / time.Second)
			c.Data["coolDownDuration"] = coolDownDurationInSecond
			c.Data["readonly"] = "readonly"
		}

		c.Data["kind"] = kind
		c.Data["name"] = name
		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Autoscaler"
	}

	guimessage.OutputMessage(c.Data)
}

func (c *EditController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	indicatorSlice := make([]Indicator, 0)
	cpu := c.GetString("cpu")
	if cpu == "on" {
		cpuAboveAllOrOneText := c.GetString("cpuAboveAllOrOne")
		var cpuAboveAllOrOne bool
		if cpuAboveAllOrOneText == "on" {
			cpuAboveAllOrOne = true
		} else {
			cpuAboveAllOrOne = false
		}
		cpuAbovePercentageOfData, _ := c.GetFloat("cpuAbovePercentageOfData")
		cpuAboveThreshold, _ := c.GetInt64("cpuAboveThreshold")
		cpuBelowAllOrOneText := c.GetString("cpuBelowAllOrOne")
		var cpuBelowAllOrOne bool
		if cpuBelowAllOrOneText == "on" {
			cpuBelowAllOrOne = true
		} else {
			cpuBelowAllOrOne = false
		}
		cpuBelowPercentageOfData, _ := c.GetFloat("cpuBelowPercentageOfData")
		cpuBelowThreshold, _ := c.GetInt64("cpuBelowThreshold")
		indicatorSlice = append(indicatorSlice, Indicator{"cpu",
			cpuAboveAllOrOne, cpuAbovePercentageOfData / 100.0, cpuAboveThreshold * 1000000,
			cpuBelowAllOrOne, cpuBelowPercentageOfData / 100.0, cpuBelowThreshold * 1000000})
	}
	memory := c.GetString("memory")
	if memory == "on" {
		memoryAboveAllOrOneText := c.GetString("memoryAboveAllOrOne")
		var memoryAboveAllOrOne bool
		if memoryAboveAllOrOneText == "on" {
			memoryAboveAllOrOne = true
		} else {
			memoryAboveAllOrOne = false
		}
		memoryAbovePercentageOfData, _ := c.GetFloat("memoryAbovePercentageOfData")
		memoryAboveThreshold, _ := c.GetInt64("memoryAboveThreshold")
		memoryBelowAllOrOneText := c.GetString("memoryBelowAllOrOne")
		var memoryBelowAllOrOne bool
		if memoryBelowAllOrOneText == "on" {
			memoryBelowAllOrOne = true
		} else {
			memoryBelowAllOrOne = false
		}
		memoryBelowPercentageOfData, _ := c.GetFloat("memoryBelowPercentageOfData")
		memoryBelowThreshold, _ := c.GetInt64("memoryBelowThreshold")
		indicatorSlice = append(indicatorSlice, Indicator{"memory",
			memoryAboveAllOrOne, memoryAbovePercentageOfData / 100.0, memoryAboveThreshold * 1024 * 1024,
			memoryBelowAllOrOne, memoryBelowPercentageOfData / 100.0, memoryBelowThreshold * 1024 * 1024})
	}

	kubernetesManagementProtocol := beego.AppConfig.String("kubernetesManagementProtocol")
	kubernetesManagementHost := beego.AppConfig.String("kubernetesManagementHost")
	kubernetesManagementPort := beego.AppConfig.String("kubernetesManagementPort")
	kubeapiHost := beego.AppConfig.String("kubeapiHost")
	kubeapiPort, _ := beego.AppConfig.Int("kubeapiPort")

	namespace, _ := c.GetSession("namespace").(string)

	kind := c.GetString("kind")
	name := c.GetString("name")
	coolDownDuration, _ := c.GetInt("coolDownDuration")
	maximumReplica, _ := c.GetInt("maximumReplica")
	minimumReplica, _ := c.GetInt("minimumReplica")

	replicationControllerAutoScaler := ReplicationControllerAutoScaler{
		true, time.Duration(coolDownDuration) * time.Second, 0, kubeapiHost, kubeapiPort, namespace, kind, name,
		maximumReplica, minimumReplica, indicatorSlice}

	url := kubernetesManagementProtocol + "://" + kubernetesManagementHost + ":" + kubernetesManagementPort +
		"/api/v1/autoscalers/"

	_, err := restclient.RequestPutWithStructure(url, replicationControllerAutoScaler, nil)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Auto scaler for " + kind + " " + name + " is edited")
	}

	c.Ctx.Redirect(302, "/gui/deploy/autoscaler/")

	guimessage.RedirectMessage(c)
}
