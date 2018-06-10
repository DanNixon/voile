package web

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func FindTitleElementInDocument(n *html.Node) (string, error) {
	if n != nil && n.FirstChild != nil {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild.Type == html.TextNode {
			return n.FirstChild.Data, nil
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			res, err := FindTitleElementInDocument(c)
			if err == nil {
				return res, nil
			}
		}
	}

	return "", errors.New("No title element found")
}

func FindTitleElement(url url.URL) (string, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	title, err := FindTitleElementInDocument(doc)
	return strings.TrimSpace(title), err
}
