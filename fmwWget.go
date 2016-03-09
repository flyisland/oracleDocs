package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type prodLink struct {
	abbr string
	name string
	href *url.URL
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
	var title = ""
	//	"../../middleware/1213/wls/index.html"
	var re = regexp.MustCompile(`middleware/\d+/(\w+)/\w+.htm`)
	var findedlink prodLink

	doc.Find(versionSelector).Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ = s.Attr("href")
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

func buildURLs(pName string) {
	mBooks := make(map[string]string)
	for _, prodLink := range prodSlices {
		// 1. get all urls wiht "cross"
		if prodLink.abbr != pName {
			continue
		}

		// 2. get into each page and find the url for "Books"
		doc, err := goquery.NewDocument(prodLink.href.String())
		if err != nil {
			log.Fatal(err)
		}

		var dirName = ""
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			if s.Text() == "Books" {
				href, _ := s.Attr("href")
				booksUrl, _ := prodLink.href.Parse(href)

				if pName == "cross" {
					dirName = "core"
				}

				strUrl := booksUrl.String()
				if _, found := mBooks[strUrl]; !found {
					findBooks(strUrl, "pdf", dirName)
					mBooks[strUrl] = strUrl
				}
			}
		})
	}
}

func findBooks(booksUrl string, args ...string) {
	var bookExt = "pdf"
	var dirRe = `(\w+)`
	if len(args) >= 1 {
		bookExt = args[0]
	}
	if len(args) >= 2 {
		if args[1] == "core" {
			dirRe = "(core)"
		}
	}

	// bookHref = "../osb/OSBAG.pdf"  -> ../dir/FILE.pdf
	var bookRe = regexp.MustCompile(dirRe + `/(\w+).` + bookExt)

	doc, err := goquery.NewDocument(booksUrl)
	if err != nil {
		log.Fatal(err)
		return
	}
	doc.Find(".booklist").Each(func(i int, s *goquery.Selection) {
		if bookHref, hrefExists := s.Find("[href$='." + bookExt + "']").Attr("href"); hrefExists {
			bookUrl, _ := doc.Url.Parse(bookHref)
			matchs := bookRe.FindStringSubmatch(bookUrl.String())
			if len(matchs) == 3 {
				// find the booktitle block
				bookTitle := s.Find(".booktitle").Text()
				bookTitle = strings.TrimSpace(bookTitle)

				for _, c := range invalidFileName {
					bookTitle = strings.Replace(bookTitle, c, "_", -1)
				}

				localBookFileName := matchs[1] + "/" + bookTitle + "." + matchs[2] + ".pdf"
				fmt.Printf("if not exist %s mkdir %s\n", matchs[1], matchs[1])
				fmt.Printf("if not exist \"%s\" wget %s -O \"%s\"\n", localBookFileName, bookUrl, localBookFileName)
			}
		}
	})
}

func readme() {
	fmt.Println("Usage: go run fmwWget.go 11g|12c PRODUCTNAME")
	fmt.Println("Build commands to download offline files for this product.")
	fmt.Println("")
	fmt.Println("PRODUCTNAME=wls     : download files for WebLogic Server.")
	fmt.Println("PRODUCTNAME=LIST    : list all products.")
}

func main() {

	if len(os.Args) != 3 {
		readme()
		os.Exit(0)
	}

	if os.Args[1] == "11g" {
		versionSelector = ".as111190" // selector for 11.1.1.9
	} else {
		versionSelector = ".rel1221" // selector for 12.2.1
	}
	listAllProducts()

	var product = os.Args[2]
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
