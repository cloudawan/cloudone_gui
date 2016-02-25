package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/volume:CreateController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/volume:CreateController"],
		beego.ControllerComments{
			"Get",
			`/createinformation/:clusterName`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/volume:CreateController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/volume:CreateController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/volume:DeleteController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/volume:DeleteController"],
		beego.ControllerComments{
			"Delete",
			`/:clusterName/:glusterfsVolume`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/volume:ListController"] = append(beego.GlobalControllerRouter["github.com/cloudawan/cloudone_gui/guirestapi/filesystem/glusterfs/volume:ListController"],
		beego.ControllerComments{
			"Get",
			`/:clusterName`,
			[]string{"get"},
			nil})

}
