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
	var baseURL = "http://docs.oracle.com/cd/E29542_01/nav/"
	var productURL = baseURL + product + ".htm"
	doc, err := goquery.NewDocument(productURL)
	if err != nil {
		log.Fatal(err)
	}

	var fullTitle = ""
	var finded = false
	var tocURL = ""
	var re = regexp.MustCompile(`([a-z]+)\.1111\/([a-z]\d+).pdf`)
	var localPdf = ""

	doc.Find(".booklist").Each(func(i int, s *goquery.Selection) {
		if pdfHref, pdfExists := s.Find("[href$='.pdf']").Attr("href"); pdfExists {
			// must be ../(dir)/(filename).pdf format
			matchs := re.FindStringSubmatch(pdfHref)
			if len(matchs) == 3 {
				fmt.Printf("if not exist %s mkdir %s\n", matchs[1], matchs[1])

				// find the booktitle block
				bb := s.Find(".booktitle")
				if len(bb.Nodes) > 0 {
					if tocURL, finded = bb.Find("a").Attr("href"); finded {
						// get the toc document and extrace the full title of this doc
						tocDoc, err := goquery.NewDocument(baseURL + tocURL)
						if err != nil {
							log.Fatal(err)
						}

						fullTitle, finded = tocDoc.Find("[name='dcterms.title']").Attr("content")
						if !finded {
							fullTitle = tocDoc.Find("title").Text()
						}
					} else {
						fullTitle = bb.Text()
					}

					fullTitle = strings.TrimSpace(fullTitle)
					fullTitle = strings.TrimLeft(fullTitle, "Fusion Middleware ")
					fullTitle = strings.TrimLeft(fullTitle, "Oracle Fusion Middleware ")
					localPdf = matchs[1] + "/" + fullTitle + "." + matchs[2] + ".pdf"
					fmt.Printf("if not exist \"%s\" wget %s%s -O \"%s\"\n", localPdf, baseURL, pdfHref, localPdf)
				}
			}
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

func listProducts() map[string]string {
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

	return m
}

func main() {
	if len(os.Args) == 1 {
		readme()
		os.Exit(0)
	}

	var product = os.Args[1]

	switch product {
	case "LIST":
		{
			m := listProducts()
			for p, s := range m {
				fmt.Printf("%s -> %s\n", p, s)
			}
		}
	case "ALL":
		{
			m := listProducts()
			for p := range m {
				buildURLs(p)
			}
		}
	default:
		buildURLs(product)
	}
}
