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

type metaType struct {
	Meta map[string]string `json:"meta"`
}

func PlaywrightScrape(c echo.Context) error {

	type hrefType struct {
		Url string `json:"hrefUrl"`
	}

	type headLinkType struct {
		Link map[string]string `json:"link"`
	}

	type Htags struct {
		TagName    string `json:"tagName"`
		TagPresent bool   `json:"tagPresent"`
		TagValue   int    `json:"tagValue"`
	}

	type imgAltType struct {
		MissingALt string `json:"alt"`
	}

	type apiResult struct {
		Title          string
		Meta           []metaType     `json:"metas"`
		Link           []headLinkType `json:"links"`
		Href           []hrefType     `json:"hrefs"`
		HTags          []Htags        `json:"htags"`
		Imgs           []imgAltType   `json:"imgs"`
		ScreenshotByte []byte         `json:"ss"`
	}

	url := c.QueryParam("url")
	fmt.Println(url)

	//if err := playwright.Install(); err != nil {
	//	log.Panic(err)
	//}

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
		return c.JSON(http.StatusRequestTimeout, "request timed out")
	}

	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
		//Timeout: playwright.Float(300),
	})

	siteHtml, _ := page.Content()

	// loading unparsed html into *io.Reader
	var parseHtml io.Reader = strings.NewReader(siteHtml)

	//parsing html
	doc, err := goquery.NewDocumentFromReader(parseHtml)
	if err != nil {
		log.Fatal(err)
	}

	//title tag
	exportTitle := doc.Find("title").Text()

	//find all meta tags
	fmt.Println("getting meta data...")
	exportMeta := []metaType{}
	doc.Find("meta").Each(func(i int, meta *goquery.Selection) {
		z := meta.Get(0)
		attributes := make(map[string]string)
		for _, att := range z.Attr {
			attributes[att.Key] = att.Val
		}
		exportMeta = append(exportMeta, metaType{attributes})
	})

	//find all links in head
	fmt.Println("getting link tags...")
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
	fmt.Println("getting hrefs...")
	exportHref := []hrefType{}
	doc.Find("a").Each(func(i int, a *goquery.Selection) {
		z, _ := a.Attr("href")
		format := hrefType{z}
		exportHref = append(exportHref, format)
	})
	//fmt.Printf("href: %v", exportHref)

	//find h* tags and return bool
	fmt.Println("getting h tags...")
	exportHTags := []Htags{}
	for i := 1; i < 7; i++ {
		conCat := "h" + strconv.Itoa(i)
		findTag := doc.Find(conCat)
		checkTag := findTag.Is(conCat)
		checkCount := findTag.Length()
		exportHTags = append(exportHTags, Htags{conCat, checkTag, checkCount})
	}
	//fmt.Printf("htags: %v", exportHTags)

	fmt.Println("getting imgs...")
	exportImgAlts := []imgAltType{}
	doc.Find("img").Each(func(i int, img *goquery.Selection) {
		alt, _ := img.Attr("alt")
		if len(alt) < 1 {
			html, _ := goquery.OuterHtml(img)
			exportImgAlts = append(exportImgAlts, imgAltType{html})
		}
	})

	screenshotByteArr, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path:    playwright.String("./temp.jpeg"),
		Quality: playwright.Int(50),
		Type:    playwright.ScreenshotTypeJpeg,
		//Timeout: playwright.Float(4000),
	})
	if err != nil {
		fmt.Println(err)
	}

	result := apiResult{
		exportTitle,
		exportMeta,
		exportLink,
		exportHref,
		exportHTags,
		exportImgAlts,
		screenshotByteArr,
	}

  if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	return c.JSON(http.StatusOK, result)
}
