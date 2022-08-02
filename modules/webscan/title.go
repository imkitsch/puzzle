package webscan

import (
	"Allin/gologger"
	"github.com/antchfx/htmlquery"
	"io"
)

func getHTTPTitle(Body io.ReadCloser) string {

	doc, err := htmlquery.Parse(Body)
	if err != nil {
		gologger.Warningf(err.Error())
	}
	title := htmlquery.FindOne(doc, `/html/head/title/text()`)
	return htmlquery.InnerText(title)
}
