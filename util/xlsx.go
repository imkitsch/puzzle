package util

import (
	"github.com/xuri/excelize/v2"
)

func XlsxInit1(output string) error {
	file := excelize.NewFile()

	//创建表
	IpIndex := file.NewSheet("IP地址")
	//CidrIndex := file.NewSheet("c段信息")
	//PortIndex := file.NewSheet("端口服务")
	//DomainFinger := file.NewSheet("域名指纹")
	//IpFinger := file.NewSheet("IP指纹")

	//初始数据
	file.SetCellValue("ip地址", "A1", "IP地址")

	file.SetCellValue("C段信息", "A1", "域名")
	file.SetCellValue("C段信息", "A2", "IP")
	file.SetCellValue("C段信息", "A3", "C段")
	file.SetCellValue("C段信息", "A4", "国家")
	file.SetCellValue("C段信息", "A5", "城市")

	file.SetCellValue("端口服务", "A1", "IP地址")
	file.SetCellValue("端口服务", "A2", "端口")
	file.SetCellValue("端口服务", "A3", "banner")

	file.SetCellValue("域名指纹", "A1", "web地址")
	file.SetCellValue("域名指纹", "A2", "IsCDN")
	file.SetCellValue("域名指纹", "A3", "状态码")
	file.SetCellValue("域名指纹", "A4", "头部信息")
	file.SetCellValue("域名指纹", "A5", "返回长度")
	file.SetCellValue("域名指纹", "A6", "标题")
	file.SetCellValue("域名指纹", "A7", "指纹")
	file.SetCellValue("域名指纹", "A8", "蜜罐")

	file.SetCellValue("IP指纹", "A1", "web地址")
	file.SetCellValue("IP指纹", "A2", "状态码")
	file.SetCellValue("IP指纹", "A3", "头部信息")
	file.SetCellValue("IP指纹", "A4", "返回长度")
	file.SetCellValue("IP指纹", "A5", "标题")
	file.SetCellValue("IP指纹", "A6", "指纹")
	file.SetCellValue("IP指纹", "A7", "蜜罐")
	file.SetActiveSheet(IpIndex)

	err := file.SaveAs(output)
	if err != nil {
		return err
	}
	return nil
}
