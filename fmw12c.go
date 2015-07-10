package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var invalidFileName = []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
var fmwBaseURL = "http://docs.oracle.com/en/middleware/"
var fmwPage = fmwBaseURL + "middleware.html"
var s12c = ".rel1213" // selector for 1213

func listProducts() map[string]string {
	doc, err := goquery.NewDocument(fmwPage)
	if err != nil {
		log.Fatal(err)
	}

	var href = ""
	var finded = false
	var shortTitle = ""
	//	"../../middleware/1213/cross/getstartedtasks.htm"
	var re = regexp.MustCompile(`middleware/\d+/(\w+)/\w+.htm`)
	var m = make(map[string]string)

	doc.Find(s12c).Find("a").Each(func(i int, s *goquery.Selection) {
		href, finded = s.Attr("href")
		shortTitle = s.Text()
		matchs := re.FindStringSubmatch(href)

		if len(matchs) > 1 {
			m[matchs[1]] = shortTitle
		}

	})

	return m
}

func buildURLs(product string) {
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
				//				bookTitle = strings.Replace(bookTitle, "Fusion Middleware ", "", 1)
				//				bookTitle = strings.Replace(bookTitle, "Oracle Fusion Middleware ", "", 1)

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

func readme() {
	fmt.Println("Usage: fmw12c PRODUCTNAME")
	fmt.Println("Build commands to download offline files for this product.")
	fmt.Println("")
	fmt.Println("PRODUCTNAME=wls     : download files for WebLogic Server.")
	fmt.Println("PRODUCTNAME=LIST    : list all products.")
	fmt.Println("PRODUCTNAME=ALL     : download files for all products!")
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
