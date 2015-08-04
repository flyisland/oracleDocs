# Download Oracle Docs

## Purpose
I'd like to download documents of Oracle's products for offline viewing, but Oracle do not provide a convenient way to do that. So I build this tool to generate `wget` commands by querying the web page.

## Procedure
1. Access the landing page like "http://docs.oracle.com/en/middleware/middleware.html"
1. For each product, go into the product's landing page, and get the link of "Books"
1. In the "Books" page, generate `wget` command for each book.

## Usage

```
Usage: go run fmwWget.go 11g|12c PRODUCTNAME
Build commands to download offline files for this product.

PRODUCTNAME=wls     : download files for WebLogic Server.
PRODUCTNAME=LIST    : list all products.
```
