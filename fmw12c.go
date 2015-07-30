package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type prodLink struct {
	abbr	string
	name	string
	href	*url.URL
}

var invalidFileName = []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
var mwURL = "http://docs.oracle.com/en/middleware/middleware.html"
var versionSelector = ".rel1213" // selector for 1213
var prodSlices = make([]prodLink, 0, 60)

func listAllProducts() {
	doc, err := goquery.NewDocument(mwURL)
	if err != nil {
		log.Fatal(err)
	}

	var href = ""
	var finded = false
	var title = ""
	//	"../../middleware/1213/wls/index.html"
	var re = regexp.MustCompile(`middleware/\d+/(\w+)/\w+.htm`)
	var findedlink	prodLink

	doc.Find(versionSelector).Find("a").Each(func(i int, s *goquery.Selection) {
		href, finded = s.Attr("href")
		title = s.Text()
		matchs := re.FindStringSubmatch(href)

		if len(matchs) > 1 {
			findedlink.abbr = matchs[1]
			findedlink.name = title
			findedlink.href, _ = doc.Url.Parse(href)
			prodSlices = append(prodSlices, findedlink)
		}
	})
}

func buildURLs(product string) {
	if (product == "cross"){
		buildCrossURLs()
		return
	}
	var baseURL = "http://docs.oracle.com/middleware/1213/" + product + "/"
	doc, err := goquery.NewDocument(baseURL + "docs.htm")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(baseURL + "docs.htm")

	// pdfHref = "../osb/OSBAG.pdf"  -> ../dir/FILE.pdf
	var re = regexp.MustCompile(`(\w+)/(\w+).pdf`)

	doc.Find(".booklist").Each(func(i int, s *goquery.Selection) {
		if pdfHref, pdfExists := s.Find("[href$='.pdf']").Attr("href"); pdfExists {
			fmt.Println(pdfHref)
			matchs := re.FindStringSubmatch(pdfHref)
			if len(matchs) == 3 {
				// find the booktitle block
				bookTitle := s.Find(".booktitle").Text()
				bookTitle = strings.TrimSpace(bookTitle)

				for _, c := range invalidFileName {
					bookTitle = strings.Replace(bookTitle, c, "_", -1)
				}

				localPdf := matchs[1] + "/" + bookTitle + "." + matchs[2] + ".pdf"
				fmt.Printf("if not exist %s mkdir %s\n", matchs[1], matchs[1])
				fmt.Printf("if not exist \"%s\" wget %s%s -O \"%s\"\n", localPdf, baseURL, pdfHref, localPdf)
			}
		}
	})
}

// Common Documents for Fusion Middleware
func buildCrossURLs() {
	// 1. get all urls wiht "cross"
	// 2. get into each page and find the url for "Books"
	// 3. search the book page find pdf under "/core/"
}

func readme() {
	fmt.Println("Usage: fmw12c PRODUCTNAME")
	fmt.Println("Build commands to download offline files for this product.")
	fmt.Println("")
	fmt.Println("PRODUCTNAME=wls     : download files for WebLogic Server.")
	fmt.Println("PRODUCTNAME=LIST    : list all products.")
}

func main() {

	if len(os.Args) == 1 {
		readme()
		os.Exit(0)
	}

	var product = os.Args[1]

	listAllProducts()
	switch product {
	case "LIST":
		{
			for i, s := range prodSlices {
				fmt.Printf("%02d. %15s -> %s\n", i, s.abbr, s.name)
			}
		}
	default:
		buildURLs(product)
	}
}
