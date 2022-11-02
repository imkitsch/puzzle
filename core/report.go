package core

import (
	"github.com/xuri/excelize/v2"
	"puzzle/gologger"
	"reflect"
)

func XlsxInit(output string) error {
	file := excelize.NewFile()

	//初始化域名表
	DomainIndex := file.NewSheet("子域名")
	DomainCategories := map[string]string{"A1": "子域名", "B1": "IsCDN", "C1": "A解析", "D1": "CNAME"}
	for k, v := range DomainCategories {
		file.SetCellValue("子域名", k, v)
	}

	//初始化ip表
	//创建表
	IpIndex := file.NewSheet("IP地址")
	//初始数据
	IpCategories := map[string]string{"A1": "IP", "B1": "City", "C1": "Country", "D1": "Source"}
	for k, v := range IpCategories {
		file.SetCellValue("IP地址", k, v)
	}
	file.SetActiveSheet(IpIndex)

	//初始化端口服务表
	PortIndex := file.NewSheet("端口服务")
	PortCategories := map[string]string{"A1": "IP", "B1": "Port", "C1": "Server", "D1": "Banner"}
	for k, v := range PortCategories {
		file.SetCellValue("端口服务", k, v)
	}
	file.SetActiveSheet(PortIndex)

	//初始化域名指纹表
	DomainFingerIndex := file.NewSheet("域名指纹")
	DomainFingerCategories := map[string]string{"A1": "Url", "B1": "IsCDN", "C1": "StatusCode", "D1": "Header", "E1": "Length", "F1": "Title", "G1": "Finger", "H1": "IsHoneypot", "I1": "priority"}
	for k, v := range DomainFingerCategories {
		file.SetCellValue("域名指纹", k, v)
	}
	file.SetActiveSheet(DomainFingerIndex)

	//初始化IP指纹表
	IpFingerIndex := file.NewSheet("IP指纹")
	IpFingerCategories := map[string]string{"A1": "Url", "B1": "StatusCode", "C1": "Header", "D1": "Length", "E1": "Title", "F1": "Finger", "G1": "IsHoneypot", "H1": "priority               0"}
	for k, v := range IpFingerCategories {
		file.SetCellValue("IP指纹", k, v)
	}
	file.SetActiveSheet(IpFingerIndex)

	file.SetActiveSheet(DomainIndex)

	file.DeleteSheet("Sheet1")
	err := file.SaveAs(output)
	if err != nil {
		return err
	}
	return nil
}

func ReportWrite(output string, sheet string, dataInterface interface{}) {
	file, err := excelize.OpenFile(output)
	if err != nil {
		gologger.Fatalf("打开文件失败:%s", err.Error())
	}

	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		gologger.Fatalf("获取写入流失败:%s", err.Error())
	}

	data := reflect.ValueOf(dataInterface)
	for i := 0; i < data.Len(); i++ {
		tmp := data.Index(i).Elem()
		row := make([]interface{}, tmp.NumField())
		for j := 0; j < tmp.NumField(); j++ {
			row[j] = excelize.Cell{Value: tmp.Field(j).Interface()}
		}
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := streamWriter.SetRow(cell, row); err != nil {
			gologger.Fatalf("写入流失败:%s", err.Error())
		}
	}

	if err := streamWriter.Flush(); err != nil {
		gologger.Fatalf("写入流失败:%s", err.Error())
	}

	file.Save()
}
