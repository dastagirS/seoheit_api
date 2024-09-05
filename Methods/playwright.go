package Methods

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func PlaywrightScrape(c echo.Context) error {

	type metaType struct {
		Meta map[string]string `json:"meta"`
	}

	type hrefType struct {
		Url string `json:"hrefUrl:"`
	}

	type headLinkType struct {
		Link map[string]string `json:"link"`
	}

	type apiResult struct {
		Meta  []metaType     `json:"metas:"`
		Link  []headLinkType `json:"links"`
		Href  []hrefType     `json:"hrefs"`
		HTags []map[string]bool     `json:"htags"`
	}

	url := c.QueryParam("url")
	fmt.Println(url)

	if err := playwright.Install(); err != nil {
		log.Panic(err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto(url); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: playwright.Float(300),
	})

	siteHtml, _ := page.Content()

	// loading unparsed html into *io.Reader
	var parseHtml io.Reader = strings.NewReader(siteHtml)

	//parsing html
	doc, err := goquery.NewDocumentFromReader(parseHtml)
	if err != nil {
		log.Fatal(err)
	}

	//find all meta tags
	exportMeta := []metaType{}
	doc.Find("meta").Each(func(i int, meta *goquery.Selection) {
		z := meta.Get(0)
		attributes := make(map[string]string)
		for _, att := range z.Attr {
			attributes[att.Key] = att.Val
		}
		exportMeta = append(exportMeta, metaType{attributes})
	})
	//fmt.Printf("meta: %v", exportMeta)4

	//find all links in head
	exportLink := []headLinkType{}
	doc.Find("link").Each(func(i int, link *goquery.Selection) {
		z := link.Get(0)
		attributes := make(map[string]string)
		for _, att := range z.Attr {
			attributes[att.Key] = att.Val
		}
		exportLink = append(exportLink, headLinkType{attributes})
	})
	//fmt.Printf("link: %v", exportLink)

	//find all href tags in body
	exportHref := []hrefType{}
	doc.Find("a").Each(func(i int, a *goquery.Selection) {
		z, _ := a.Attr("href")
		format := hrefType{z}
		exportHref = append(exportHref, format)
	})
	//fmt.Printf("href: %v", exportHref)

	//find h* tags and return bool
	var exportHTags []map[string]bool
	for i := 1; i < 7; i++ {
		conCat := "h" + strconv.Itoa(i)
		checkTag := doc.Find(conCat).Is(conCat)
		tag := make(map[string]bool)
		tag[conCat] = checkTag
		exportHTags = append(exportHTags, tag)
	}
	//fmt.Printf("htags: %v", exportHTags)

	result := apiResult{
		exportMeta,
		exportLink,
		exportHref,
		exportHTags,
	}

	fmt.Println(result)

	return c.JSON(http.StatusOK, result)
}
