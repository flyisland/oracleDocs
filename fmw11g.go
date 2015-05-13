package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func buildURLs() {
	baseURL := "http://docs.oracle.com/cd/E29542_01/nav/"

	doc, err := goquery.NewDocument(baseURL + "wlst.htm")
	if err != nil {
		log.Fatal(err)
	}

	var fullTitle = ""
	var finded = false
	var tocURL = ""
	var re = regexp.MustCompile(`([a-z]+)\.1111\/([a-z]\d+).pdf`) // (dir)(filename)

	doc.Find(".booklist").Each(func(i int, s *goquery.Selection) {
		// find the booktitle block
		tocURL, finded = s.Find(".booktitle > a").Attr("href")
		if !finded {
			log.Fatal("Can not found booktitle!")
		}

		pdfHref, pdfExists := s.Find("[href$='.pdf']").Attr("href")

		if pdfExists {
			// get the toc document and extrace the full title of this doc
			tocDoc, err := goquery.NewDocument(baseURL + tocURL)
			if err != nil {
				log.Fatal(err)
			}
			fullTitle, finded = tocDoc.Find("[name='dcterms.title']").Attr("content")
			if !finded {
				fullTitle = tocDoc.Find("title").Text()
			}

			fullTitle = strings.TrimLeft(fullTitle, "Fusion Middleware ")
			fullTitle = strings.TrimLeft(fullTitle, "Oracle Fusion Middleware ")
			matchs := re.FindStringSubmatch(pdfHref)
			localPdf := matchs[1] + "/" + strings.TrimSpace(fullTitle) + "." + matchs[2] + ".pdf"

			fmt.Printf("if not exist %s mkdir %s\n", matchs[1], matchs[1])
			fmt.Printf("if not exist \"%s\" wget %s%s -O \"%s\"\n", localPdf, baseURL, pdfHref, localPdf)
		}

	})
}

func main() {
	buildURLs()
}
