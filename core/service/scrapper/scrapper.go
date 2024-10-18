package scrapperService

import (
	"api_journal/core/server/shared"
	"api_journal/core/service/postagensService"
	"api_journal/core/util"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func RunScrapper() {

	fmt.Println("Lets scrapp")
	currentTitleReport := ""
	currentTitleNews := ""

	reportPageCollector := colly.NewCollector()
	newsPageCollector := colly.NewCollector()

	newsPageCollector.OnHTML("section#content", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, elem *colly.HTMLElement) {
			link := elem.Attr("href")
			if !strings.Contains(link, "https://www.ufca.edu.br/noticias/page/") && link != "https://www.ufca.edu.br" {
				newsPageCollector.Visit(link)
			}
		})
	})

	newsPageCollector.OnHTML("article#content", func(e *colly.HTMLElement) {

		e.ForEach("h1", func(x int, elem *colly.HTMLElement) {
			currentTitleNews = elem.Text
		})
		updateDataDiv := e.DOM.Find("p.content-update-data")
		htmlAfter, _ := updateDataDiv.NextAll().Html()
		novaData := util.StringToDate(updateDataDiv.Text())
		postagensService.InsertPostFromUfca(shared.DB, "ufca-noticia", currentTitleNews, htmlAfter, []string{""}, novaData)
	})

	reportPageCollector.OnHTML("article#content", func(e *colly.HTMLElement) {

		e.ForEach("h1", func(x int, elem *colly.HTMLElement) {
			currentTitleReport = elem.Text
		})
		updateDataDiv := e.DOM.Find("p.content-update-data")
		htmlAfter, _ := e.DOM.Find("div.content.one.column.row").Html()
		novaData := util.StringToDate(updateDataDiv.Text())
		postagensService.InsertPostFromUfca(shared.DB, "ufca-reportagem", currentTitleReport, htmlAfter, []string{""}, novaData)
	})

	reportPageCollector.OnHTML("section.twelve", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(x int, elem *colly.HTMLElement) {
			link := elem.Attr("href")
			contains := strings.Contains(link, "https://www.ufca.edu.br/informes/page")
			if !contains && link != "https://www.ufca.edu.br" {
				reportPageCollector.Visit(link)
			}
		})
	})

	reportPageCollector.Visit("https://www.ufca.edu.br/informes/")
	newsPageCollector.Visit("https://www.ufca.edu.br/noticias/")
	reportPageCollector.Wait()
	newsPageCollector.Wait()

}
