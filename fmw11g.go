package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func buildURLs(product string) {
	baseURL := "http://docs.oracle.com/cd/E29542_01/nav/"

	doc, err := goquery.NewDocument(baseURL + product + ".htm")
	if err != nil {
		log.Fatal(err)
	}

	var fullTitle = ""
	var finded = false
	var tocURL = ""
	var re = regexp.MustCompile(`([a-z]+)\.1111\/([a-z]\d+).pdf`) // (dir)(filename)

	doc.Find(".booklist").Each(func(i int, s *goquery.Selection) {
		pdfHref, pdfExists := s.Find("[href$='.pdf']").Attr("href")

		if pdfExists {
			// find the booktitle block
			tocURL, finded = s.Find(".booktitle > a").Attr("href")
			if !finded {
				log.Fatal("Can not found toc url of " + pdfHref)
			}
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

func readme() {
	fmt.Println("Usage: fmw11g PRODUCTNAME")
	fmt.Println("Build commands to download offline files for this product.")
	fmt.Println("")
	fmt.Println("PRODUCTNAME=wls     : download files for WebLogic Server.")
	fmt.Println("PRODUCTNAME=LIST    : list all products.")
	fmt.Println("PRODUCTNAME=ALL     : download files for all products!")
}

func listProducts() {
	// Oracle Fusion Middleware Online Documentation Library 11g Release 1 (11.1.1.8)
	baseURL := "http://docs.oracle.com/cd/E29542_01/index.htm"
	doc, err := goquery.NewDocument(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	var href = ""
	var finded = false
	var shortTitle = ""
	var re = regexp.MustCompile(`(\w+)\.`)
	var m = make(map[string]string)

	doc.Find("[style='text-decoration:none']").Each(func(i int, s *goquery.Selection) {
		href, finded = s.Attr("href")
		shortTitle = s.Text()
		matchs := re.FindStringSubmatch(href)

		if len(matchs) > 1 {
			m[matchs[1]] = shortTitle
		}

	})

	fmt.Println(m)
}

func main() {
	if len(os.Args) == 1 {
		readme()
		os.Exit(0)
	}

	var product = os.Args[1]

	switch product {
	case "LIST":
		listProducts()
	case "ALL":
	default:
		buildURLs(product)
	}
}
