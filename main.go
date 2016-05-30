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

package main

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	_ "github.com/cloudawan/cloudone_gui/docs" // Import document generation
	restapiidentity "github.com/cloudawan/cloudone_gui/restapi/v1/identity"
	_ "github.com/cloudawan/cloudone_gui/routers"
)

func main() {
	beego.InsertFilter("/gui/*", beego.BeforeRouter, identity.FilterUser)
	beego.InsertFilter("/api/v1/*", beego.BeforeRouter, restapiidentity.FilterToken)
	beego.InsertFilter("/guirestapi/v1/*", beego.BeforeRouter, identity.FilterUser)

	beego.AppConfigPath = "/etc/cloudone_gui/app.conf"

	beego.Run()
}
