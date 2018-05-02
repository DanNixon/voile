package web

import (
	"errors"
	"net/http"

	"golang.org/x/net/html"
)

func FinditleElementInDocument(n *html.Node) (string, error) {
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild.Type == html.TextNode {
		return n.FirstChild.Data, nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res, err := FinditleElementInDocument(c)
		if err == nil {
			return res, nil
		}
	}

	return "", errors.New("No title element found")
}

func FindTitleElement(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	return FinditleElementInDocument(doc)
}
