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

package container

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/dashboard"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"time"
)

type ReplicationControllerMetricList struct {
	ErrorSlice                       []string
	ReplicationControllerMetricSlice []ReplicationControllerMetric
}

type ReplicationControllerMetric struct {
	Namespace                 string
	ReplicationControllerName string
	ValidPodSlice             []bool
	PodMetricSlice            []PodMetric
	Size                      int
}

type PodMetric struct {
	KubeletHost          string
	Namespace            string
	PodName              string
	ValidContainerSlice  []bool
	ContainerMetricSlice []ContainerMetric
}

type ContainerMetric struct {
	ContainerName                     string
	CpuUsageTotalSlice                []int64
	MemoryUsageSlice                  []int64
	DiskIOServiceBytesStatsTotalSlice []int64
	DiskIOServicedStatsTotalSlice     []int64
	NetworkRXBytesSlice               []int64
	NetworkTXBytesSlice               []int64
	NetworkRXPacketsSlice             []int64
	NetworkTXPacketsSlice             []int64
}

type IndexController struct {
	beego.Controller
}

const (
	allKeyword = "All"
)

func (c *IndexController) Get() {
	c.TplName = "monitor/container/index.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	cloudoneGUIProtocol := beego.AppConfig.String("cloudoneGUIProtocol")
	cloudoneGUIHost, cloudoneGUIPort := dashboard.GetServerHostAndPortFromUserRequest(c.Ctx.Input)

	namespaces, _ := c.GetSession("namespace").(string)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/replicationcontrollers/" + namespaces

	jsonMapSlice := make([]interface{}, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &jsonMapSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(guimessagedisplay.GetErrorMessage(err))
	} else {
		replicationControllerNameSlice := make([]string, 0)
		replicationControllerNameSlice = append(replicationControllerNameSlice, allKeyword)

		for _, jsonMap := range jsonMapSlice {
			name, _ := jsonMap.(map[string]interface{})["Name"].(string)
			if name != "" {
				replicationControllerNameSlice = append(replicationControllerNameSlice, name)
			}
		}

		c.Data["cloudoneGUIProtocol"] = cloudoneGUIProtocol
		c.Data["cloudoneGUIHost"] = cloudoneGUIHost
		c.Data["cloudoneGUIPort"] = cloudoneGUIPort
		c.Data["replicationControllerNameSlice"] = replicationControllerNameSlice
	}

	guimessage.OutputMessage(c.Data)
}

type DataController struct {
	beego.Controller
}

func (c *DataController) Get() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	namespaces, _ := c.GetSession("namespace").(string)

	replicationControllerName := c.GetString("replicationController")

	replicationControllerMetricSlice := make([]ReplicationControllerMetric, 0)
	replicationControllerMetricAmount := 0
	if replicationControllerName != "" && replicationControllerName != allKeyword {
		replicationControllerMetric := ReplicationControllerMetric{}
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/replicationcontrollermetrics/" + namespaces + "/" + replicationControllerName

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestGetWithStructure(url, &replicationControllerMetric, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			errorJsonMap := make(map[string]interface{})
			errorJsonMap["error"] = err.Error()
			c.Data["json"] = errorJsonMap
			c.ServeJSON()
			return
		}
		replicationControllerMetricSlice = append(replicationControllerMetricSlice, replicationControllerMetric)
		replicationControllerMetricAmount = 1
	} else {
		url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
			"/api/v1/replicationcontrollermetrics/" + namespaces

		replicationControllerMetricList := ReplicationControllerMetricList{}

		tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

		_, err := restclient.RequestGetWithStructure(url, &replicationControllerMetricList, tokenHeaderMap)

		if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
			return
		}

		if err != nil {
			// Error
			errorJsonMap := make(map[string]interface{})
			errorJsonMap["error"] = err.Error()
			c.Data["json"] = errorJsonMap
			c.ServeJSON()
			return
		}
		replicationControllerMetricSlice = replicationControllerMetricList.ReplicationControllerMetricSlice
		replicationControllerMetricAmount = len(replicationControllerMetricSlice)
	}

	// CPU usage total
	cpuUsageTotalJsonMap := make(map[string]interface{})
	cpuUsageTotalDifferenceAmountMaximum := 0
	// Memory usage
	memoryUsageJsonMap := make(map[string]interface{})
	memoryUsageAmountMaximum := 0
	// Disk I/O byte
	diskIOServiceBytesStatsJsonMap := make(map[string]interface{})
	diskIOServiceBytesStatsDifferenceAmountMaximum := 0
	// Disk I/O count
	diskIOServicedStatsJsonMap := make(map[string]interface{})
	diskIOServicedStatsDifferenceAmountMaximum := 0
	// Network RX Byte
	networkRXBytesJsonMap := make(map[string]interface{})
	networkRXBytesDifferenceAmountMaximum := 0
	// Network TX Byte
	networkTXBytesJsonMap := make(map[string]interface{})
	networkTXBytesDifferenceAmountMaximum := 0
	// Network RX Packets
	networkRXPacketsJsonMap := make(map[string]interface{})
	networkRXPacketsDifferenceAmountMaximum := 0
	// Network TX Packets
	networkTXPacketsJsonMap := make(map[string]interface{})
	networkTXPacketsDifferenceAmountMaximum := 0

	if replicationControllerMetricAmount > 0 {
		// CPU usage total
		cpuUsageTotalJsonMap["series"] = make([]interface{}, 0)

		// Memory usage
		memoryUsageJsonMap["series"] = make([]interface{}, 0)

		// Disk I/O byte
		diskIOServiceBytesStatsJsonMap["series"] = make([]interface{}, 0)

		// Disk I/O count
		diskIOServicedStatsJsonMap["series"] = make([]interface{}, 0)

		// Network RX Byte
		networkRXBytesJsonMap["series"] = make([]interface{}, 0)

		// Network TX Byte
		networkTXBytesJsonMap["series"] = make([]interface{}, 0)

		// Network RX Packets
		networkRXPacketsJsonMap["series"] = make([]interface{}, 0)

		// Network TX Packets
		networkTXPacketsJsonMap["series"] = make([]interface{}, 0)

		for _, replicationControllerMetric := range replicationControllerMetricSlice {
			podMetricSlice := replicationControllerMetric.PodMetricSlice
			for _, podMetric := range podMetricSlice {
				containerMetricSlice := podMetric.ContainerMetricSlice
				for _, containerMetric := range containerMetricSlice {
					// CPU usage total
					cpuUsageTotalDifferenceAmount := len(containerMetric.CpuUsageTotalSlice) - 1
					if cpuUsageTotalDifferenceAmount > cpuUsageTotalDifferenceAmountMaximum {
						cpuUsageTotalDifferenceAmountMaximum = cpuUsageTotalDifferenceAmount
					}
					cpuUsageTotalJsonMapSeries := make(map[string]interface{})
					cpuUsageTotalJsonMapSeries["name"] = replicationControllerMetric.ReplicationControllerName + "/" + podMetric.PodName + "/" + containerMetric.ContainerName
					cpuUsageTotalJsonMapSeries["data"] = make([]int, cpuUsageTotalDifferenceAmount)
					for j := 0; j < cpuUsageTotalDifferenceAmount; j++ {
						cpuUsageTotalJsonMapSeries["data"].([]int)[j] = int(containerMetric.CpuUsageTotalSlice[j+1]-containerMetric.CpuUsageTotalSlice[j]) / 1000000
					}
					cpuUsageTotalJsonMap["series"] = append(cpuUsageTotalJsonMap["series"].([]interface{}), cpuUsageTotalJsonMapSeries)
					// Memory usage
					memoryUsageAmount := len(containerMetric.MemoryUsageSlice)
					if memoryUsageAmount > memoryUsageAmountMaximum {
						memoryUsageAmountMaximum = memoryUsageAmount
					}
					memoryUsageJsonMapSeries := make(map[string]interface{})
					memoryUsageJsonMapSeries["name"] = replicationControllerMetric.ReplicationControllerName + "/" + podMetric.PodName + "/" + containerMetric.ContainerName
					memoryUsageJsonMapSeries["data"] = make([]int, memoryUsageAmount)
					for j := 0; j < memoryUsageAmount; j++ {
						memoryUsageJsonMapSeries["data"].([]int)[j] = int(containerMetric.MemoryUsageSlice[j]) / (1024 * 1024)
					}
					memoryUsageJsonMap["series"] = append(memoryUsageJsonMap["series"].([]interface{}), memoryUsageJsonMapSeries)
					// Disk I/O byte
					diskIOServiceBytesStatsDifferenceAmount := len(containerMetric.DiskIOServiceBytesStatsTotalSlice) - 1
					if diskIOServiceBytesStatsDifferenceAmount > diskIOServiceBytesStatsDifferenceAmountMaximum {
						diskIOServiceBytesStatsDifferenceAmountMaximum = diskIOServiceBytesStatsDifferenceAmount
					}
					diskIOServiceBytesStatsJsonMapSeries := make(map[string]interface{})
					diskIOServiceBytesStatsJsonMapSeries["name"] = replicationControllerMetric.ReplicationControllerName + "/" + podMetric.PodName + "/" + containerMetric.ContainerName
					diskIOServiceBytesStatsJsonMapSeries["data"] = make([]int, diskIOServiceBytesStatsDifferenceAmount)
					for j := 0; j < diskIOServiceBytesStatsDifferenceAmount; j++ {
						diskIOServiceBytesStatsJsonMapSeries["data"].([]int)[j] = int(containerMetric.DiskIOServiceBytesStatsTotalSlice[j+1] - containerMetric.DiskIOServiceBytesStatsTotalSlice[j])
					}
					diskIOServiceBytesStatsJsonMap["series"] = append(diskIOServiceBytesStatsJsonMap["series"].([]interface{}), diskIOServiceBytesStatsJsonMapSeries)
					// Disk I/O count
					diskIOServicedStatsDifferenceAmount := len(containerMetric.DiskIOServicedStatsTotalSlice) - 1
					if diskIOServicedStatsDifferenceAmount > diskIOServicedStatsDifferenceAmountMaximum {
						diskIOServicedStatsDifferenceAmountMaximum = diskIOServicedStatsDifferenceAmount
					}
					diskIOServicedStatsJsonMapSeries := make(map[string]interface{})
					diskIOServicedStatsJsonMapSeries["name"] = replicationControllerMetric.ReplicationControllerName + "/" + podMetric.PodName + "/" + containerMetric.ContainerName
					diskIOServicedStatsJsonMapSeries["data"] = make([]int, diskIOServicedStatsDifferenceAmount)
					for j := 0; j < diskIOServicedStatsDifferenceAmount; j++ {
						diskIOServicedStatsJsonMapSeries["data"].([]int)[j] = int(containerMetric.DiskIOServicedStatsTotalSlice[j+1] - containerMetric.DiskIOServicedStatsTotalSlice[j])
					}
					diskIOServicedStatsJsonMap["series"] = append(diskIOServicedStatsJsonMap["series"].([]interface{}), diskIOServicedStatsJsonMapSeries)
					// Network RX Byte
					networkRXBytesDifferenceAmount := len(containerMetric.NetworkRXBytesSlice) - 1
					if networkRXBytesDifferenceAmount > networkRXBytesDifferenceAmountMaximum {
						networkRXBytesDifferenceAmountMaximum = networkRXBytesDifferenceAmount
					}
					networkRXBytesJsonMapSeries := make(map[string]interface{})
					networkRXBytesJsonMapSeries["name"] = replicationControllerMetric.ReplicationControllerName + "/" + podMetric.PodName + "/" + containerMetric.ContainerName
					networkRXBytesJsonMapSeries["data"] = make([]int, networkRXBytesDifferenceAmount)
					for j := 0; j < networkRXBytesDifferenceAmount; j++ {
						networkRXBytesJsonMapSeries["data"].([]int)[j] = int(containerMetric.NetworkRXBytesSlice[j+1] - containerMetric.NetworkRXBytesSlice[j])
					}
					networkRXBytesJsonMap["series"] = append(networkRXBytesJsonMap["series"].([]interface{}), networkRXBytesJsonMapSeries)
					// Network TX Byte
					networkTXBytesDifferenceAmount := len(containerMetric.NetworkTXBytesSlice) - 1
					if networkTXBytesDifferenceAmount > networkTXBytesDifferenceAmountMaximum {
						networkTXBytesDifferenceAmountMaximum = networkTXBytesDifferenceAmount
					}
					networkTXBytesJsonMapSeries := make(map[string]interface{})
					networkTXBytesJsonMapSeries["name"] = replicationControllerMetric.ReplicationControllerName + "/" + podMetric.PodName + "/" + containerMetric.ContainerName
					networkTXBytesJsonMapSeries["data"] = make([]int, networkTXBytesDifferenceAmount)
					for j := 0; j < networkTXBytesDifferenceAmount; j++ {
						networkTXBytesJsonMapSeries["data"].([]int)[j] = int(containerMetric.NetworkTXBytesSlice[j+1] - containerMetric.NetworkTXBytesSlice[j])
					}
					networkTXBytesJsonMap["series"] = append(networkTXBytesJsonMap["series"].([]interface{}), networkTXBytesJsonMapSeries)
					// Network RX Packet
					networkRXPacketsDifferenceAmount := len(containerMetric.NetworkRXPacketsSlice) - 1
					if networkRXPacketsDifferenceAmount > networkRXPacketsDifferenceAmountMaximum {
						networkRXPacketsDifferenceAmountMaximum = networkRXPacketsDifferenceAmount
					}
					networkRXPacketsJsonMapSeries := make(map[string]interface{})
					networkRXPacketsJsonMapSeries["name"] = replicationControllerMetric.ReplicationControllerName + "/" + podMetric.PodName + "/" + containerMetric.ContainerName
					networkRXPacketsJsonMapSeries["data"] = make([]int, networkRXPacketsDifferenceAmount)
					for j := 0; j < networkRXPacketsDifferenceAmount; j++ {
						networkRXPacketsJsonMapSeries["data"].([]int)[j] = int(containerMetric.NetworkRXPacketsSlice[j+1] - containerMetric.NetworkRXPacketsSlice[j])
					}
					networkRXPacketsJsonMap["series"] = append(networkRXPacketsJsonMap["series"].([]interface{}), networkRXPacketsJsonMapSeries)
					// Network TX Packet
					networkTXPacketsDifferenceAmount := len(containerMetric.NetworkTXPacketsSlice) - 1
					if networkTXPacketsDifferenceAmount > networkTXPacketsDifferenceAmountMaximum {
						networkTXPacketsDifferenceAmountMaximum = networkTXPacketsDifferenceAmount
					}
					networkTXPacketsJsonMapSeries := make(map[string]interface{})
					networkTXPacketsJsonMapSeries["name"] = replicationControllerMetric.ReplicationControllerName + "/" + podMetric.PodName + "/" + containerMetric.ContainerName
					networkTXPacketsJsonMapSeries["data"] = make([]int, networkTXPacketsDifferenceAmount)
					for j := 0; j < networkTXPacketsDifferenceAmount; j++ {
						networkTXPacketsJsonMapSeries["data"].([]int)[j] = int(containerMetric.NetworkTXPacketsSlice[j+1] - containerMetric.NetworkTXPacketsSlice[j])
					}
					networkTXPacketsJsonMap["series"] = append(networkTXPacketsJsonMap["series"].([]interface{}), networkTXPacketsJsonMapSeries)
				}
			}
		}
		// CPU usage total
		cpuUsageTotalJsonMap["labels"] = make([]int, cpuUsageTotalDifferenceAmountMaximum)
		for i := 0; i < cpuUsageTotalDifferenceAmountMaximum; i++ {
			cpuUsageTotalJsonMap["labels"].([]int)[i] = -1*cpuUsageTotalDifferenceAmountMaximum + i + 1
		}
		// Memory usage
		memoryUsageJsonMap["labels"] = make([]int, memoryUsageAmountMaximum)
		for i := 0; i < memoryUsageAmountMaximum; i++ {
			memoryUsageJsonMap["labels"].([]int)[i] = -1*memoryUsageAmountMaximum + i + 1
		}
		// Disk I/O byte
		diskIOServiceBytesStatsJsonMap["labels"] = make([]int, diskIOServiceBytesStatsDifferenceAmountMaximum)
		for i := 0; i < diskIOServiceBytesStatsDifferenceAmountMaximum; i++ {
			diskIOServiceBytesStatsJsonMap["labels"].([]int)[i] = -1*diskIOServiceBytesStatsDifferenceAmountMaximum + i + 1
		}
		// Disk I/O count
		diskIOServicedStatsJsonMap["labels"] = make([]int, diskIOServicedStatsDifferenceAmountMaximum)
		for i := 0; i < diskIOServicedStatsDifferenceAmountMaximum; i++ {
			diskIOServicedStatsJsonMap["labels"].([]int)[i] = -1*diskIOServicedStatsDifferenceAmountMaximum + i + 1
		}
		// Network RX Byte
		networkRXBytesJsonMap["labels"] = make([]int, networkRXBytesDifferenceAmountMaximum)
		for i := 0; i < networkRXBytesDifferenceAmountMaximum; i++ {
			networkRXBytesJsonMap["labels"].([]int)[i] = -1*networkRXBytesDifferenceAmountMaximum + i + 1
		}
		// Network TX Byte
		networkTXBytesJsonMap["labels"] = make([]int, networkTXBytesDifferenceAmountMaximum)
		for i := 0; i < networkTXBytesDifferenceAmountMaximum; i++ {
			networkTXBytesJsonMap["labels"].([]int)[i] = -1*networkTXBytesDifferenceAmountMaximum + i + 1
		}
		// Network RX Packet
		networkRXPacketsJsonMap["labels"] = make([]int, networkRXPacketsDifferenceAmountMaximum)
		for i := 0; i < networkRXPacketsDifferenceAmountMaximum; i++ {
			networkRXPacketsJsonMap["labels"].([]int)[i] = -1*networkRXPacketsDifferenceAmountMaximum + i + 1
		}
		// Network RX Packet
		networkTXPacketsJsonMap["labels"] = make([]int, networkTXPacketsDifferenceAmountMaximum)
		for i := 0; i < networkTXPacketsDifferenceAmountMaximum; i++ {
			networkTXPacketsJsonMap["labels"].([]int)[i] = -1*networkTXPacketsDifferenceAmountMaximum + i + 1
		}
	}

	// Convert
	current := time.Now()
	// CPU usage total
	convertedCpuUsageTotalJsonMap := make(map[string]interface{})
	convertedCpuUsageTotalJsonMap["metadata"] = make(map[string]interface{})
	convertedCpuUsageTotalJsonMap["metadata"].(map[string]interface{})["title"] = "CPU (ms/1s)"
	convertedCpuUsageTotalJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	convertedCpuUsageTotalJsonMap["data"] = make([]map[string]interface{}, cpuUsageTotalDifferenceAmountMaximum)
	for i := 0; i < cpuUsageTotalDifferenceAmountMaximum; i++ {
		convertedCpuUsageTotalJsonMap["data"].([]map[string]interface{})[i] = make(map[string]interface{})
		convertedCpuUsageTotalJsonMap["data"].([]map[string]interface{})[i]["y"] = make([]int, 0)
		convertedCpuUsageTotalJsonMap["data"].([]map[string]interface{})[i]["x"] = current.Add(
			time.Duration(time.Second * time.Duration(-1*(cpuUsageTotalDifferenceAmountMaximum-i)))).Format("2006-01-02 15:04:05")
	}
	seriesJsonMap, _ := cpuUsageTotalJsonMap["series"].([]interface{})
	for _, line := range seriesJsonMap {
		convertedCpuUsageTotalJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
			convertedCpuUsageTotalJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
			line.(map[string]interface{})["name"].(string))
		for i := 0; i < len(line.(map[string]interface{})["data"].([]int)); i++ {
			convertedCpuUsageTotalJsonMap["data"].([]map[string]interface{})[i]["y"] = append(
				convertedCpuUsageTotalJsonMap["data"].([]map[string]interface{})[i]["y"].([]int),
				line.(map[string]interface{})["data"].([]int)[i])
		}
	}
	// Memory usage
	convertedMemoryUsageJsonMap := make(map[string]interface{})
	convertedMemoryUsageJsonMap["metadata"] = make(map[string]interface{})
	convertedMemoryUsageJsonMap["metadata"].(map[string]interface{})["title"] = "Memory(MB)"
	convertedMemoryUsageJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	convertedMemoryUsageJsonMap["data"] = make([]map[string]interface{}, memoryUsageAmountMaximum)
	for i := 0; i < memoryUsageAmountMaximum; i++ {
		convertedMemoryUsageJsonMap["data"].([]map[string]interface{})[i] = make(map[string]interface{})
		convertedMemoryUsageJsonMap["data"].([]map[string]interface{})[i]["y"] = make([]int, 0)
		convertedMemoryUsageJsonMap["data"].([]map[string]interface{})[i]["x"] = current.Add(
			time.Duration(time.Second * time.Duration(-1*(memoryUsageAmountMaximum-i)))).Format("2006-01-02 15:04:05")
	}
	seriesJsonMap, _ = memoryUsageJsonMap["series"].([]interface{})
	for _, line := range seriesJsonMap {
		convertedMemoryUsageJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
			convertedMemoryUsageJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
			line.(map[string]interface{})["name"].(string))
		for i := 0; i < len(line.(map[string]interface{})["data"].([]int)); i++ {
			convertedMemoryUsageJsonMap["data"].([]map[string]interface{})[i]["y"] = append(
				convertedMemoryUsageJsonMap["data"].([]map[string]interface{})[i]["y"].([]int),
				line.(map[string]interface{})["data"].([]int)[i])
		}
	}
	// Disk I/O byte
	convertedDiskIOServiceBytesStatsJsonMap := make(map[string]interface{})
	convertedDiskIOServiceBytesStatsJsonMap["metadata"] = make(map[string]interface{})
	convertedDiskIOServiceBytesStatsJsonMap["metadata"].(map[string]interface{})["title"] = "Disk I/O (Byte/s)"
	convertedDiskIOServiceBytesStatsJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	convertedDiskIOServiceBytesStatsJsonMap["data"] = make([]map[string]interface{}, diskIOServiceBytesStatsDifferenceAmountMaximum)
	for i := 0; i < diskIOServiceBytesStatsDifferenceAmountMaximum; i++ {
		convertedDiskIOServiceBytesStatsJsonMap["data"].([]map[string]interface{})[i] = make(map[string]interface{})
		convertedDiskIOServiceBytesStatsJsonMap["data"].([]map[string]interface{})[i]["y"] = make([]int, 0)
		convertedDiskIOServiceBytesStatsJsonMap["data"].([]map[string]interface{})[i]["x"] = current.Add(
			time.Duration(time.Second * time.Duration(-1*(diskIOServiceBytesStatsDifferenceAmountMaximum-i)))).Format("2006-01-02 15:04:05")
	}
	seriesJsonMap, _ = diskIOServiceBytesStatsJsonMap["series"].([]interface{})
	for _, line := range seriesJsonMap {
		convertedDiskIOServiceBytesStatsJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
			convertedDiskIOServiceBytesStatsJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
			line.(map[string]interface{})["name"].(string))
		for i := 0; i < len(line.(map[string]interface{})["data"].([]int)); i++ {
			convertedDiskIOServiceBytesStatsJsonMap["data"].([]map[string]interface{})[i]["y"] = append(
				convertedDiskIOServiceBytesStatsJsonMap["data"].([]map[string]interface{})[i]["y"].([]int),
				line.(map[string]interface{})["data"].([]int)[i])
		}
	}
	// Disk I/O count
	convertedDiskIOServicedStatsJsonMap := make(map[string]interface{})
	convertedDiskIOServicedStatsJsonMap["metadata"] = make(map[string]interface{})
	convertedDiskIOServicedStatsJsonMap["metadata"].(map[string]interface{})["title"] = "Disk I/O (count/s)"
	convertedDiskIOServicedStatsJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	convertedDiskIOServicedStatsJsonMap["data"] = make([]map[string]interface{}, diskIOServicedStatsDifferenceAmountMaximum)
	for i := 0; i < diskIOServicedStatsDifferenceAmountMaximum; i++ {
		convertedDiskIOServicedStatsJsonMap["data"].([]map[string]interface{})[i] = make(map[string]interface{})
		convertedDiskIOServicedStatsJsonMap["data"].([]map[string]interface{})[i]["y"] = make([]int, 0)
		convertedDiskIOServicedStatsJsonMap["data"].([]map[string]interface{})[i]["x"] = current.Add(
			time.Duration(time.Second * time.Duration(-1*(diskIOServicedStatsDifferenceAmountMaximum-i)))).Format("2006-01-02 15:04:05")
	}
	seriesJsonMap, _ = diskIOServicedStatsJsonMap["series"].([]interface{})
	for _, line := range seriesJsonMap {
		convertedDiskIOServicedStatsJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
			convertedDiskIOServicedStatsJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
			line.(map[string]interface{})["name"].(string))
		for i := 0; i < len(line.(map[string]interface{})["data"].([]int)); i++ {
			convertedDiskIOServicedStatsJsonMap["data"].([]map[string]interface{})[i]["y"] = append(
				convertedDiskIOServicedStatsJsonMap["data"].([]map[string]interface{})[i]["y"].([]int),
				line.(map[string]interface{})["data"].([]int)[i])
		}
	}
	// Network RX Byte
	convertedNetworkRXBytesJsonMap := make(map[string]interface{})
	convertedNetworkRXBytesJsonMap["metadata"] = make(map[string]interface{})
	convertedNetworkRXBytesJsonMap["metadata"].(map[string]interface{})["title"] = "Network RX (Bytes/s)"
	convertedNetworkRXBytesJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	convertedNetworkRXBytesJsonMap["data"] = make([]map[string]interface{}, networkRXBytesDifferenceAmountMaximum)
	for i := 0; i < networkRXBytesDifferenceAmountMaximum; i++ {
		convertedNetworkRXBytesJsonMap["data"].([]map[string]interface{})[i] = make(map[string]interface{})
		convertedNetworkRXBytesJsonMap["data"].([]map[string]interface{})[i]["y"] = make([]int, 0)
		convertedNetworkRXBytesJsonMap["data"].([]map[string]interface{})[i]["x"] = current.Add(
			time.Duration(time.Second * time.Duration(-1*(networkRXBytesDifferenceAmountMaximum-i)))).Format("2006-01-02 15:04:05")
	}
	seriesJsonMap, _ = networkRXBytesJsonMap["series"].([]interface{})
	for _, line := range seriesJsonMap {
		convertedNetworkRXBytesJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
			convertedNetworkRXBytesJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
			line.(map[string]interface{})["name"].(string))
		for i := 0; i < len(line.(map[string]interface{})["data"].([]int)); i++ {
			convertedNetworkRXBytesJsonMap["data"].([]map[string]interface{})[i]["y"] = append(
				convertedNetworkRXBytesJsonMap["data"].([]map[string]interface{})[i]["y"].([]int),
				line.(map[string]interface{})["data"].([]int)[i])
		}
	}
	// Network TX Byte
	convertedNetworkTXBytesJsonMap := make(map[string]interface{})
	convertedNetworkTXBytesJsonMap["metadata"] = make(map[string]interface{})
	convertedNetworkTXBytesJsonMap["metadata"].(map[string]interface{})["title"] = "Network TX (Bytes/s)"
	convertedNetworkTXBytesJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	convertedNetworkTXBytesJsonMap["data"] = make([]map[string]interface{}, networkTXBytesDifferenceAmountMaximum)
	for i := 0; i < networkTXBytesDifferenceAmountMaximum; i++ {
		convertedNetworkTXBytesJsonMap["data"].([]map[string]interface{})[i] = make(map[string]interface{})
		convertedNetworkTXBytesJsonMap["data"].([]map[string]interface{})[i]["y"] = make([]int, 0)
		convertedNetworkTXBytesJsonMap["data"].([]map[string]interface{})[i]["x"] = current.Add(
			time.Duration(time.Second * time.Duration(-1*(networkTXBytesDifferenceAmountMaximum-i)))).Format("2006-01-02 15:04:05")
	}
	seriesJsonMap, _ = networkTXBytesJsonMap["series"].([]interface{})
	for _, line := range seriesJsonMap {
		convertedNetworkTXBytesJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
			convertedNetworkTXBytesJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
			line.(map[string]interface{})["name"].(string))
		for i := 0; i < len(line.(map[string]interface{})["data"].([]int)); i++ {
			convertedNetworkTXBytesJsonMap["data"].([]map[string]interface{})[i]["y"] = append(
				convertedNetworkTXBytesJsonMap["data"].([]map[string]interface{})[i]["y"].([]int),
				line.(map[string]interface{})["data"].([]int)[i])
		}
	}
	// Network RX Packet
	convertedNetworkRXPacketsJsonMap := make(map[string]interface{})
	convertedNetworkRXPacketsJsonMap["metadata"] = make(map[string]interface{})
	convertedNetworkRXPacketsJsonMap["metadata"].(map[string]interface{})["title"] = "Network RX (packet/s)"
	convertedNetworkRXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	convertedNetworkRXPacketsJsonMap["data"] = make([]map[string]interface{}, networkRXPacketsDifferenceAmountMaximum)
	for i := 0; i < networkRXPacketsDifferenceAmountMaximum; i++ {
		convertedNetworkRXPacketsJsonMap["data"].([]map[string]interface{})[i] = make(map[string]interface{})
		convertedNetworkRXPacketsJsonMap["data"].([]map[string]interface{})[i]["y"] = make([]int, 0)
		convertedNetworkRXPacketsJsonMap["data"].([]map[string]interface{})[i]["x"] = current.Add(
			time.Duration(time.Second * time.Duration(-1*(networkRXPacketsDifferenceAmountMaximum-i)))).Format("2006-01-02 15:04:05")
	}
	seriesJsonMap, _ = networkRXPacketsJsonMap["series"].([]interface{})
	for _, line := range seriesJsonMap {
		convertedNetworkRXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
			convertedNetworkRXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
			line.(map[string]interface{})["name"].(string))
		for i := 0; i < len(line.(map[string]interface{})["data"].([]int)); i++ {
			convertedNetworkRXPacketsJsonMap["data"].([]map[string]interface{})[i]["y"] = append(
				convertedNetworkRXPacketsJsonMap["data"].([]map[string]interface{})[i]["y"].([]int),
				line.(map[string]interface{})["data"].([]int)[i])
		}
	}
	// Network TX Packet
	convertedNetworkTXPacketsJsonMap := make(map[string]interface{})
	convertedNetworkTXPacketsJsonMap["metadata"] = make(map[string]interface{})
	convertedNetworkTXPacketsJsonMap["metadata"].(map[string]interface{})["title"] = "Network TX (packet/s)"
	convertedNetworkTXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"] = make([]string, 0)
	convertedNetworkTXPacketsJsonMap["data"] = make([]map[string]interface{}, networkTXPacketsDifferenceAmountMaximum)
	for i := 0; i < networkTXPacketsDifferenceAmountMaximum; i++ {
		convertedNetworkTXPacketsJsonMap["data"].([]map[string]interface{})[i] = make(map[string]interface{})
		convertedNetworkTXPacketsJsonMap["data"].([]map[string]interface{})[i]["y"] = make([]int, 0)
		convertedNetworkTXPacketsJsonMap["data"].([]map[string]interface{})[i]["x"] = current.Add(
			time.Duration(time.Second * time.Duration(-1*(networkTXPacketsDifferenceAmountMaximum-i)))).Format("2006-01-02 15:04:05")
	}
	seriesJsonMap, _ = networkTXPacketsJsonMap["series"].([]interface{})
	for _, line := range seriesJsonMap {
		convertedNetworkTXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
			convertedNetworkTXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
			line.(map[string]interface{})["name"].(string))
		for i := 0; i < len(line.(map[string]interface{})["data"].([]int)); i++ {
			convertedNetworkTXPacketsJsonMap["data"].([]map[string]interface{})[i]["y"] = append(
				convertedNetworkTXPacketsJsonMap["data"].([]map[string]interface{})[i]["y"].([]int),
				line.(map[string]interface{})["data"].([]int)[i])
		}
	}

	chartSlice := make(map[string]interface{})
	// CPU usage total
	chartSlice["cpuUsageTotal"] = convertedCpuUsageTotalJsonMap
	// Memory usage
	chartSlice["memoryUsage"] = convertedMemoryUsageJsonMap
	// Disk I/O byte
	chartSlice["diskIOServiceBytesStats"] = convertedDiskIOServiceBytesStatsJsonMap
	// Disk I/O count
	chartSlice["diskIOServicedStats"] = convertedDiskIOServicedStatsJsonMap
	// Network RX Byte
	chartSlice["networkRXBytes"] = convertedNetworkRXBytesJsonMap
	// Network TX Byte
	chartSlice["networkTXBytes"] = convertedNetworkTXBytesJsonMap
	// Network RX Packet
	chartSlice["networkRXPackets"] = convertedNetworkRXPacketsJsonMap
	// Network TX Packet
	chartSlice["networkTXPackets"] = convertedNetworkTXPacketsJsonMap

	c.Data["json"] = chartSlice

	c.ServeJSON()
}
