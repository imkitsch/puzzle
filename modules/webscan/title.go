package webscan

import (
	"github.com/antchfx/htmlquery"
	"io"
	"puzzle/gologger"
)

func getHTTPTitle(Body io.ReadCloser) string {

	doc, err := htmlquery.Parse(Body)
	if err != nil {
		gologger.Warningf(err.Error())
	}
	title := htmlquery.FindOne(doc, `/html/head/title/text()`)
	return htmlquery.InnerText(title)
}
