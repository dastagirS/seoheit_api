package Methods

import (
	"fmt"
)

type metaAttrType struct {
	metaAttr  string
	metaValue []string
}

func MetaChecker(siteMetaArr []metaType) error {

	metaChecklist := []metaAttrType{
		{
			metaAttr:  "name",
			metaValue: []string{"author", "contact", "twitter:creator", "twitter:card", "twitter:url", "twitter:title", "twitter:description", "twitter:image:alt", "usp", "robots", "keywords", "description", "msapplication-TileImage"},
		},
		{
			metaAttr:  "property",
			metaValue: []string{"og:url", "og:type", "og:image", "og:site_name", "og:description", "og:locale", "og:image:alt", "og:title", "twitter:url", "twitter:image"},
		},
	}

	_ = metaChecklist


	//for _, metaArrays := range metaChecklist {
		//if metaArrays.metaAttr == "name" {
			for _, metaArr := range siteMetaArr {
		for _, meta := range metaArr.Meta {
			fmt.Println(meta)
		}
			}

	//	}
	//}

	return nil

}
