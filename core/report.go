package core

import "github.com/xuri/excelize/v2"

func XlsxInit(output string) error {
	file := excelize.NewFile()

	//创建表
	IpIndex := file.NewSheet("IP地址")
	//初始数据
	categories := map[string]string{"A1": "IP地址", "B1": "City", "C1": "Country", "D1": "Cidr"}
	for k, v := range categories {
		file.SetCellValue("ip地址", k, v)
	}
	file.DeleteSheet("Sheet1")
	file.SetActiveSheet(IpIndex)
	err := file.SaveAs(output)
	if err != nil {
		return err
	}
	return nil
}

func ReportIp(output string) {

}

func ReportCidr(output string) {

}

func ReportPort(output string) {

}

func ReportDomainFinger(output string) {

}

func ReportIpFinger(output string) {

}
