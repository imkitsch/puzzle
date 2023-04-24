package webscan

import (
	"github.com/antchfx/htmlquery"
	"io"
)

func getHTTPTitle(Body *io.ReadCloser) string {

	doc, err := htmlquery.Parse(*Body)
	if err != nil {
		return ""
	}
	title, err := htmlquery.Query(doc, `/html/head/title/text()`)
	if err != nil || title == nil {
		return ""
	}
	return htmlquery.InnerText(title)
}
