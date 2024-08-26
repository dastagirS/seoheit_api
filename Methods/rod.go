package Methods

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/labstack/echo/v4"
)

func GetSite(c echo.Context) error {
	url := c.Param("*")

	page := rod.New().MustConnect().MustPage(url)

  wait  := page.WaitRequestIdle(300 *time.Millisecond, nil, nil, nil)
  wait()
	defer page.Close()

	siteHtml, _ := page.HTML()
	var parseHtml io.Reader = strings.NewReader(siteHtml)


  type siteInfo struct{

  }

	doc, err := goquery.NewDocumentFromReader(parseHtml)
	if err != nil {
		log.Fatal(err)
	}

  type metaType struct{
    tagName string
    tagValue string
  }

	doc.Find("meta").Each(func(i int, meta *goquery.Selection) {
		z := meta.Get(0)
		for _, att := range z.Attr {
			fmt.Println(att.Key, att.Val)
		}
	})

	doc.Find("link").Each(func(i int, link *goquery.Selection) {
		z := link.Get(0)
		for _, att := range z.Attr {
			fmt.Println(att.Key, att.Val)
		}
	})

	/*doc.Find("h1").Each(func(i int, hTags *goquery.Selection) {
		z := hTags.Get(0)
		fmt.Println(z.Data)
	})*/



	type Htag struct {
		TagName   string
		IsPresent bool
	}

	hTags :=[]Htag{}
	for i := 1; i < 7; i++ {
		conCat := "h" + strconv.Itoa(i)
		checkTag := doc.Find(conCat).Is(conCat)
    //tagIngo := conCat + strconv.FormatBool(checkTag)
    format := Htag{
      TagName: conCat,
      IsPresent: checkTag,
    }
		hTags = append(hTags, format)
	}

  fmt.Println(hTags)

	return c.JSON(http.StatusOK, siteHtml)
}
