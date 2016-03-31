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

package node

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"time"
)

type NodeMetric struct {
	Valid                             bool
	KubeletHost                       string
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

func (c *IndexController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	c.TplName = "monitor/node/index.html"

	cloudoneGUIProtocol := beego.AppConfig.String("cloudoneGUIProtocol")
	cloudoneGUIHost := c.Ctx.Input.Host()
	cloudoneGUIPort := c.Ctx.Input.Port()

	c.Data["cloudoneGUIProtocol"] = cloudoneGUIProtocol
	c.Data["cloudoneGUIHost"] = cloudoneGUIHost
	c.Data["cloudoneGUIPort"] = cloudoneGUIPort

	guimessage.OutputMessage(c.Data)
}

type DataController struct {
	beego.Controller
}

func (c *DataController) Get() {

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.ServeJSON()
		return
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/nodemetrics/?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	nodeMetricSlice := make([]NodeMetric, 0)

	_, err = restclient.RequestGetWithStructure(url, &nodeMetricSlice)

	if err != nil {
		// Error
		errorJsonMap := make(map[string]interface{})
		errorJsonMap["error"] = err.Error()
		c.Data["json"] = errorJsonMap
		c.ServeJSON()
		return
	}

	nodeAmount := 0
	for _, nodeMetric := range nodeMetricSlice {
		if nodeMetric.Valid {
			nodeAmount++
		}
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

	if nodeAmount > 0 {
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

		for _, nodeMetric := range nodeMetricSlice {
			if nodeMetric.Valid {
				// CPU usage total
				cpuUsageTotalDifferenceAmount := len(nodeMetric.CpuUsageTotalSlice) - 1
				if cpuUsageTotalDifferenceAmount > cpuUsageTotalDifferenceAmountMaximum {
					cpuUsageTotalDifferenceAmountMaximum = cpuUsageTotalDifferenceAmount
				}
				cpuUsageTotalJsonMapSeries := make(map[string]interface{})
				cpuUsageTotalJsonMapSeries["name"] = nodeMetric.KubeletHost
				cpuUsageTotalJsonMapSeries["data"] = make([]int, cpuUsageTotalDifferenceAmount)
				for j := 0; j < cpuUsageTotalDifferenceAmount; j++ {
					cpuUsageTotalJsonMapSeries["data"].([]int)[j] = int(nodeMetric.CpuUsageTotalSlice[j+1]-nodeMetric.CpuUsageTotalSlice[j]) / 1000000
				}
				cpuUsageTotalJsonMap["series"] = append(cpuUsageTotalJsonMap["series"].([]interface{}), cpuUsageTotalJsonMapSeries)
				// Memory usage
				memoryUsageAmount := len(nodeMetric.MemoryUsageSlice)
				if memoryUsageAmount > memoryUsageAmountMaximum {
					memoryUsageAmountMaximum = memoryUsageAmount
				}
				memoryUsageJsonMapSeries := make(map[string]interface{})
				memoryUsageJsonMapSeries["name"] = nodeMetric.KubeletHost
				memoryUsageJsonMapSeries["data"] = make([]int, memoryUsageAmount)
				for j := 0; j < memoryUsageAmount; j++ {
					memoryUsageJsonMapSeries["data"].([]int)[j] = int(nodeMetric.MemoryUsageSlice[j]) / (1024 * 1024)
				}
				memoryUsageJsonMap["series"] = append(memoryUsageJsonMap["series"].([]interface{}), memoryUsageJsonMapSeries)
				// Disk I/O byte
				diskIOServiceBytesStatsDifferenceAmount := len(nodeMetric.DiskIOServiceBytesStatsTotalSlice) - 1
				if diskIOServiceBytesStatsDifferenceAmount > diskIOServiceBytesStatsDifferenceAmountMaximum {
					diskIOServiceBytesStatsDifferenceAmountMaximum = diskIOServiceBytesStatsDifferenceAmount
				}
				diskIOServiceBytesStatsJsonMapSeries := make(map[string]interface{})
				diskIOServiceBytesStatsJsonMapSeries["name"] = nodeMetric.KubeletHost
				diskIOServiceBytesStatsJsonMapSeries["data"] = make([]int, diskIOServiceBytesStatsDifferenceAmount)
				for j := 0; j < diskIOServiceBytesStatsDifferenceAmount; j++ {
					diskIOServiceBytesStatsJsonMapSeries["data"].([]int)[j] = int(nodeMetric.DiskIOServiceBytesStatsTotalSlice[j+1] - nodeMetric.DiskIOServiceBytesStatsTotalSlice[j])
				}
				diskIOServiceBytesStatsJsonMap["series"] = append(diskIOServiceBytesStatsJsonMap["series"].([]interface{}), diskIOServiceBytesStatsJsonMapSeries)
				// Disk I/O count
				diskIOServicedStatsDifferenceAmount := len(nodeMetric.DiskIOServicedStatsTotalSlice) - 1
				if diskIOServicedStatsDifferenceAmount > diskIOServicedStatsDifferenceAmountMaximum {
					diskIOServicedStatsDifferenceAmountMaximum = diskIOServicedStatsDifferenceAmount
				}
				diskIOServicedStatsJsonMapSeries := make(map[string]interface{})
				diskIOServicedStatsJsonMapSeries["name"] = nodeMetric.KubeletHost
				diskIOServicedStatsJsonMapSeries["data"] = make([]int, diskIOServicedStatsDifferenceAmount)
				for j := 0; j < diskIOServicedStatsDifferenceAmount; j++ {
					diskIOServicedStatsJsonMapSeries["data"].([]int)[j] = int(nodeMetric.DiskIOServicedStatsTotalSlice[j+1] - nodeMetric.DiskIOServicedStatsTotalSlice[j])
				}
				diskIOServicedStatsJsonMap["series"] = append(diskIOServicedStatsJsonMap["series"].([]interface{}), diskIOServicedStatsJsonMapSeries)
				// Network RX Byte
				networkRXBytesDifferenceAmount := len(nodeMetric.NetworkRXBytesSlice) - 1
				if networkRXBytesDifferenceAmount > networkRXBytesDifferenceAmountMaximum {
					networkRXBytesDifferenceAmountMaximum = networkRXBytesDifferenceAmount
				}
				networkRXBytesJsonMapSeries := make(map[string]interface{})
				networkRXBytesJsonMapSeries["name"] = nodeMetric.KubeletHost
				networkRXBytesJsonMapSeries["data"] = make([]int, networkRXBytesDifferenceAmount)
				for j := 0; j < networkRXBytesDifferenceAmount; j++ {
					networkRXBytesJsonMapSeries["data"].([]int)[j] = int(nodeMetric.NetworkRXBytesSlice[j+1] - nodeMetric.NetworkRXBytesSlice[j])
				}
				networkRXBytesJsonMap["series"] = append(networkRXBytesJsonMap["series"].([]interface{}), networkRXBytesJsonMapSeries)
				// Network TX Byte
				networkTXBytesDifferenceAmount := len(nodeMetric.NetworkTXBytesSlice) - 1
				if networkTXBytesDifferenceAmount > networkTXBytesDifferenceAmountMaximum {
					networkTXBytesDifferenceAmountMaximum = networkTXBytesDifferenceAmount
				}
				networkTXBytesJsonMapSeries := make(map[string]interface{})
				networkTXBytesJsonMapSeries["name"] = nodeMetric.KubeletHost
				networkTXBytesJsonMapSeries["data"] = make([]int, networkTXBytesDifferenceAmount)
				for j := 0; j < networkTXBytesDifferenceAmount; j++ {
					networkTXBytesJsonMapSeries["data"].([]int)[j] = int(nodeMetric.NetworkTXBytesSlice[j+1] - nodeMetric.NetworkTXBytesSlice[j])
				}
				networkTXBytesJsonMap["series"] = append(networkTXBytesJsonMap["series"].([]interface{}), networkTXBytesJsonMapSeries)
				// Network RX Packet
				networkRXPacketsDifferenceAmount := len(nodeMetric.NetworkRXPacketsSlice) - 1
				if networkRXPacketsDifferenceAmount > networkRXPacketsDifferenceAmountMaximum {
					networkRXPacketsDifferenceAmountMaximum = networkRXPacketsDifferenceAmount
				}
				networkRXPacketsJsonMapSeries := make(map[string]interface{})
				networkRXPacketsJsonMapSeries["name"] = nodeMetric.KubeletHost
				networkRXPacketsJsonMapSeries["data"] = make([]int, networkRXPacketsDifferenceAmount)
				for j := 0; j < networkRXPacketsDifferenceAmount; j++ {
					networkRXPacketsJsonMapSeries["data"].([]int)[j] = int(nodeMetric.NetworkRXPacketsSlice[j+1] - nodeMetric.NetworkRXPacketsSlice[j])
				}
				networkRXPacketsJsonMap["series"] = append(networkRXPacketsJsonMap["series"].([]interface{}), networkRXPacketsJsonMapSeries)
				// Network TX Packet
				networkTXPacketsDifferenceAmount := len(nodeMetric.NetworkTXPacketsSlice) - 1
				if networkTXPacketsDifferenceAmount > networkTXPacketsDifferenceAmountMaximum {
					networkTXPacketsDifferenceAmountMaximum = networkTXPacketsDifferenceAmount
				}
				networkTXPacketsJsonMapSeries := make(map[string]interface{})
				networkTXPacketsJsonMapSeries["name"] = nodeMetric.KubeletHost
				networkTXPacketsJsonMapSeries["data"] = make([]int, networkTXPacketsDifferenceAmount)
				for j := 0; j < networkTXPacketsDifferenceAmount; j++ {
					networkTXPacketsJsonMapSeries["data"].([]int)[j] = int(nodeMetric.NetworkTXPacketsSlice[j+1] - nodeMetric.NetworkTXPacketsSlice[j])
				}
				networkTXPacketsJsonMap["series"] = append(networkTXPacketsJsonMap["series"].([]interface{}), networkTXPacketsJsonMapSeries)
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
	for _, line := range cpuUsageTotalJsonMap["series"].([]interface{}) {
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
	for _, line := range memoryUsageJsonMap["series"].([]interface{}) {
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
	for _, line := range diskIOServiceBytesStatsJsonMap["series"].([]interface{}) {
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
	for _, line := range diskIOServicedStatsJsonMap["series"].([]interface{}) {
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
	for _, line := range networkRXBytesJsonMap["series"].([]interface{}) {
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
	for _, line := range networkTXBytesJsonMap["series"].([]interface{}) {
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
	for _, line := range networkRXPacketsJsonMap["series"].([]interface{}) {
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
	for _, line := range networkTXPacketsJsonMap["series"].([]interface{}) {
		convertedNetworkTXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"] = append(
			convertedNetworkTXPacketsJsonMap["metadata"].(map[string]interface{})["lineName"].([]string),
			line.(map[string]interface{})["name"].(string))
		for i := 0; i < len(line.(map[string]interface{})["data"].([]int)); i++ {
			convertedNetworkTXPacketsJsonMap["data"].([]map[string]interface{})[i]["y"] = append(
				convertedNetworkTXPacketsJsonMap["data"].([]map[string]interface{})[i]["y"].([]int),
				line.(map[string]interface{})["data"].([]int)[i])
		}
	}

	chartJsonMap := make(map[string]interface{})
	// Memory usage
	chartJsonMap["cpuUsageTotal"] = convertedCpuUsageTotalJsonMap
	// Memory usage
	chartJsonMap["memoryUsage"] = convertedMemoryUsageJsonMap
	// Disk I/O byte
	chartJsonMap["diskIOServiceBytesStats"] = convertedDiskIOServiceBytesStatsJsonMap
	// Disk I/O count
	chartJsonMap["diskIOServicedStats"] = convertedDiskIOServicedStatsJsonMap
	// Network RX Byte
	chartJsonMap["networkRXBytes"] = convertedNetworkRXBytesJsonMap
	// Network TX Byte
	chartJsonMap["networkTXBytes"] = convertedNetworkTXBytesJsonMap
	// Network RX Packet
	chartJsonMap["networkRXPackets"] = convertedNetworkRXPacketsJsonMap
	// Network TX Packet
	chartJsonMap["networkTXPackets"] = convertedNetworkTXPacketsJsonMap

	c.Data["json"] = chartJsonMap

	c.ServeJSON()
}
