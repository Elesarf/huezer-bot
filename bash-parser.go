package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetCitate ...
func GetCitate() string {

	var answer string
	bashURL := "http://bash.im/random"

	doc, err := goquery.NewDocument(bashURL)
	_check(err)

	doc.Find(".quote__body").Each(func(i int, s *goquery.Selection) {
		str, _ := s.Html()
		result := strings.ReplaceAll(str, "<br/>", "\n")
		result = strings.ReplaceAll(result, "&#34;", "\"")
		result = strings.ReplaceAll(result, "&lt;", "<")
		result = strings.ReplaceAll(result, "&gt;", ">")
		answer = result
		fmt.Println(result)
		fmt.Println("have")
	})

	return answer
}
