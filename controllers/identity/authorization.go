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

package identity

import (
	"bytes"
	"github.com/cloudawan/cloudone_utility/rbac"
)

func GetConponentName() string {
	return componentName
}

func SetPriviledgeHiddenTag(data map[interface{}]interface{}, tagName string, user *rbac.User, method string, path string) {
	if user.HasPermission(componentName, method, path) {
		data[tagName] = "<div class='btn-group'>"
	} else {
		data[tagName] = "<div hidden>"
	}
}

func GetLayoutMenu(user *rbac.User) string {
	if user == nil {
		return ""
	}

	buffer := bytes.Buffer{}
	buffer.WriteByte('\n')

	// Dashboard
	if user.HasChildPermission(componentName, "GET", "/gui/dashboard") {
		buffer.WriteString("					<li class=''><a href='/gui/dashboard/topology'>Dashboard</a></li>\n")
	}

	// Repository
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/repository") {
		buffer.WriteString("					<li class='dropdown'>\n")
		buffer.WriteString("						<a href='#' class='dropdown-toggle' data-toggle='dropdown' role='button' aria-expanded='false'>Repository<span class='caret'></span></a>\n")
		buffer.WriteString("						<ul class='dropdown-menu' role='menu'>\n")
	}
	// Child
	if user.HasPermission(componentName, "GET", "/gui/repository/imageinformation/list") {
		buffer.WriteString("							<li><a href='/gui/repository/imageinformation/list'>Applications</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/repository/thirdparty/list") {
		buffer.WriteString("							<li><a href='/gui/repository/thirdparty/list'>Third-party Services</a></li>\n")
	}
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/repository") {
		buffer.WriteString("						</ul>\n")
		buffer.WriteString("					</li>\n")
	}

	// Deploy
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/deploy") {
		buffer.WriteString("					<li class='dropdown'>\n")
		buffer.WriteString("						<a href='#' class='dropdown-toggle' data-toggle='dropdown' role='button' aria-expanded='false'>Deploy<span class='caret'></span></a>\n")
		buffer.WriteString("						<ul class='dropdown-menu' role='menu'>\n")
	}
	// Child
	if user.HasPermission(componentName, "GET", "/gui/deploy/deploy/list") {
		buffer.WriteString("							<li><a href='/gui/deploy/deploy/list'>Applications</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/deploy/deploybluegreen/list") {
		buffer.WriteString("							<li><a href='/gui/deploy/deploybluegreen/list'>Blue Green Deployments</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/deploy/autoscaler/list") {
		buffer.WriteString("							<li><a href='/gui/deploy/autoscaler/list'>Autoscalers</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/deploy/deployclusterapplication/list") {
		buffer.WriteString("							<li><a href='/gui/deploy/deployclusterapplication/list'>Third-party Services</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/deploy/clone/topology") {
		buffer.WriteString("							<li><a href='/gui/deploy/clone/select'>Clone Topology</a></li>\n")
	}
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/deploy") {
		buffer.WriteString("						</ul>\n")
		buffer.WriteString("					</li>\n")
	}

	// Inventory
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/inventory") {
		buffer.WriteString("					<li class='dropdown'>\n")
		buffer.WriteString("						<a href='#' class='dropdown-toggle' data-toggle='dropdown' role='button' aria-expanded='false'>Inventory<span class='caret'></span></a>\n")
		buffer.WriteString("						<ul class='dropdown-menu' role='menu'>\n")
	}
	// Child
	if user.HasPermission(componentName, "GET", "/gui/inventory/replicationcontroller/list") {
		buffer.WriteString("							<li><a href='/gui/inventory/replicationcontroller/list'>Replications</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/inventory/service/list") {
		buffer.WriteString("							<li><a href='/gui/inventory/service/list'>Services</a></li>\n")
	}
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/inventory") {
		buffer.WriteString("						</ul>\n")
		buffer.WriteString("					</li>\n")
	}

	// Filesystem
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/filesystem") {
		buffer.WriteString("					<li class='dropdown'>\n")
		buffer.WriteString("						<a href='#' class='dropdown-toggle' data-toggle='dropdown' role='button' aria-expanded='false'>Filesystem<span class='caret'></span></a>\n")
		buffer.WriteString("						<ul class='dropdown-menu' role='menu'>\n")
	}
	// Child
	if user.HasPermission(componentName, "GET", "/gui/filesystem/glusterfs/cluster/list") {
		buffer.WriteString("							<li><a href='/gui/filesystem/glusterfs/cluster/list'>Glusterfs</a></li>\n")
	}
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/filesystem") {
		buffer.WriteString("						</ul>\n")
		buffer.WriteString("					</li>\n")
	}

	// Monitor
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/monitor") {
		buffer.WriteString("					<li class='dropdown'>\n")
		buffer.WriteString("						<a href='#' class='dropdown-toggle' data-toggle='dropdown' role='button' aria-expanded='false'>Monitor<span class='caret'></span></a>\n")
		buffer.WriteString("						<ul class='dropdown-menu' role='menu'>\n")
	}
	// Child
	if user.HasPermission(componentName, "GET", "/gui/monitor/node/list") {
		buffer.WriteString("							<li><a href='/gui/monitor/node'>Nodes</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/monitor/container/list") {
		buffer.WriteString("							<li><a href='/gui/monitor/container'>Containers</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/monitor/historicalcontainer/list") {
		buffer.WriteString("							<li><a href='/gui/monitor/historicalcontainer'>Historical Containers</a></li>\n")
	}
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/monitor") {
		buffer.WriteString("						</ul>\n")
		buffer.WriteString("					</li>\n")
	}

	// Event
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/event") {
		buffer.WriteString("					<li class='dropdown'>\n")
		buffer.WriteString("						<a href='#' class='dropdown-toggle' data-toggle='dropdown' role='button' aria-expanded='false'>Event<span class='caret'></span></a>\n")
		buffer.WriteString("						<ul class='dropdown-menu' role='menu'>\n")
	}
	// Child
	if user.HasPermission(componentName, "GET", "/gui/event/audit/list") {
		buffer.WriteString("							<li><a href='/gui/event/audit/list'>Audit Logs</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/event/kubernetes/list") {
		buffer.WriteString("							<li><a href='/gui/event/kubernetes/list'>Kubernetes Events</a></li>\n")
	}
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/event") {
		buffer.WriteString("						</ul>\n")
		buffer.WriteString("					</li>\n")
	}

	// Notification
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/notification") {
		buffer.WriteString("					<li class='dropdown'>\n")
		buffer.WriteString("						<a href='#' class='dropdown-toggle' data-toggle='dropdown' role='button' aria-expanded='false'>Notification<span class='caret'></span></a>\n")
		buffer.WriteString("						<ul class='dropdown-menu' role='menu'>\n")
	}
	// Child
	if user.HasPermission(componentName, "GET", "/gui/notification/notifier/list") {
		buffer.WriteString("							<li><a href='/gui/notification/notifier/list'>Notifiers</a></li>\n")
	}
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/notification") {
		buffer.WriteString("						</ul>\n")
		buffer.WriteString("					</li>\n")
	}

	// System
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/system") {
		buffer.WriteString("					<li class='dropdown'>\n")
		buffer.WriteString("						<a href='#' class='dropdown-toggle' data-toggle='dropdown' role='button' aria-expanded='false'>System<span class='caret'></span></a>\n")
		buffer.WriteString("						<ul class='dropdown-menu' role='menu'>\n")
	}
	// Child
	if user.HasPermission(componentName, "GET", "/gui/system/about") {
		buffer.WriteString("							<li><a href='/gui/system/about'>About</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/system/namespace/list") {
		buffer.WriteString("							<li><a href='/gui/system/namespace/list'>Namepaces</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/system/notification/emailserver/list") {
		buffer.WriteString("							<li><a href='/gui/system/notification/emailserver/list'>Notification</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/system/host/credential/list") {
		buffer.WriteString("							<li><a href='/gui/system/host/credential/list'>Host Credential</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/system/rbac/user/list") {
		buffer.WriteString("							<li><a href='/gui/system/rbac/user/list'>RBAC</a></li>\n")
	}
	if user.HasPermission(componentName, "GET", "/gui/system/upgrade/list") {
		buffer.WriteString("							<li><a href='/gui/system/upgrade'>Upgrade</a></li>\n")
	}
	// Parent
	if user.HasChildPermission(componentName, "GET", "/gui/system") {
		buffer.WriteString("						</ul>\n")
		buffer.WriteString("					</li>\n")
	}

	return buffer.String()
}

func GetDashboardTabMenu(user *rbac.User, activeTab string) string {
	if user == nil {
		return ""
	}

	buffer := bytes.Buffer{}
	buffer.WriteByte('\n')

	if user.HasPermission(componentName, "GET", "/gui/dashboard/topology") {
		if activeTab == "topology" {
			buffer.WriteString("			<li role='presentation' class='active'><a href='#' role='tab' >Topology</a></li>\n")
		} else {
			buffer.WriteString("			<li role='presentation'><a href='/gui/dashboard/topology/' role='tab' >Topology</a></li>\n")
		}
	}
	if user.HasPermission(componentName, "GET", "/gui/dashboard/healthcheck") {
		if activeTab == "healthcheck" {
			buffer.WriteString("			<li role='presentation' class='active'><a href='#' role='tab' >Health Check</a></li>\n")
		} else {
			buffer.WriteString("			<li role='presentation'><a href='/gui/dashboard/healthcheck/list' role='tab' >Health Check</a></li>\n")
		}
	}
	if user.HasPermission(componentName, "GET", "/gui/dashboard/bluegreen") {
		if activeTab == "bluegreen" {
			buffer.WriteString("			<li role='presentation' class='active'><a href='#' role='tab' >BlueGreen</a></li>\n")
		} else {
			buffer.WriteString("			<li role='presentation'><a href='/gui/dashboard/bluegreen/' role='tab' >BlueGreen</a></li>\n")
		}
	}
	if user.HasPermission(componentName, "GET", "/gui/dashboard/appservice") {
		if activeTab == "appservice" {
			buffer.WriteString("			<li role='presentation' class='active'><a href='#' role='tab' >AppService</a></li>\n")
		} else {
			buffer.WriteString("			<li role='presentation'><a href='/gui/dashboard/appservice/' role='tab' >AppService</a></li>\n")
		}
	}
	if user.HasPermission(componentName, "GET", "/gui/dashboard/deploy") {
		if activeTab == "deploy" {
			buffer.WriteString("			<li role='presentation' class='active'><a href='#' role='tab' >Deployment</a></li>\n")
		} else {
			buffer.WriteString("			<li role='presentation'><a href='/gui/dashboard/deploy/' role='tab' >Deployment</a></li>\n")
		}
	}

	return buffer.String()
}

func GetSystemNotificationTabMenu(user *rbac.User, activeTab string) string {
	if user == nil {
		return ""
	}

	buffer := bytes.Buffer{}
	buffer.WriteByte('\n')

	if user.HasPermission(componentName, "GET", "/gui/system/notification/emailserver/list") {
		if activeTab == "emailserver" {
			buffer.WriteString("			<li role='presentation' class='active'><a href='#' role='tab' >Email Server</a></li>\n")
		} else {
			buffer.WriteString("			<li role='presentation'><a href='/gui/system/notification/emailserver/list' role='tab' >Email Server</a></li>\n")
		}
	}
	if user.HasPermission(componentName, "GET", "/gui/system/notification/sms/list") {
		if activeTab == "sms" {
			buffer.WriteString("			<li role='presentation' class='active'><a href='#' role='tab' >SMS</a></li>\n")
		} else {
			buffer.WriteString("			<li role='presentation'><a href='/gui/system/notification/sms/list' role='tab' >SMS</a></li>\n")
		}
	}
	return buffer.String()
}

func GetSystemRBACTabMenu(user *rbac.User, activeTab string) string {
	if user == nil {
		return ""
	}

	buffer := bytes.Buffer{}
	buffer.WriteByte('\n')

	if user.HasPermission(componentName, "GET", "/gui/system/rbac/user/list") {
		if activeTab == "user" {
			buffer.WriteString("			<li role='presentation' class='active'><a href='#' role='tab' >User</a></li>\n")
		} else {
			buffer.WriteString("			<li role='presentation'><a href='/gui/system/rbac/user/list' role='tab' >User</a></li>\n")
		}
	}
	if user.HasPermission(componentName, "GET", "/gui/system/rbac/role/list") {
		if activeTab == "role" {
			buffer.WriteString("			<li role='presentation' class='active'><a href='#' role='tab' >Role</a></li>\n")
		} else {
			buffer.WriteString("			<li role='presentation'><a href='/gui/system/rbac/role/list' role='tab' >Role</a></li>\n")
		}
	}
	return buffer.String()
}
