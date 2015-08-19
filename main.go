package main

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/identity"
	restapiidentity "github.com/cloudawan/kubernetes_management_gui/restapi/v1/identity"
	_ "github.com/cloudawan/kubernetes_management_gui/routers"
)

func main() {
	beego.InsertFilter("/gui/*", beego.BeforeRouter, identity.FilterUser)
	beego.InsertFilter("/api/v1/*", beego.BeforeRouter, restapiidentity.FilterToken)

	beego.Run()
}
