package api

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

type youdao struct{}

const youdaoURL = "https://www.youdao.com/w/eng/%s"

func (dict youdao) Search(word string) {
	block, err := dict.getContent(word)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	pronunciation := dict.parsePronunciation(block)
	explanation := dict.parseExplanation(block)
	dict.display(word, pronunciation, explanation)
}

func (dict youdao) getContent(word string) (*goquery.Selection, error) {
	url := fmt.Sprintf(youdaoURL, word)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	block := doc.Find("#phrsListTab")
	return block, nil
}

func (dict youdao) parsePronunciation(block *goquery.Selection) []string {
	pronunciation := make([]string, 0, 2)
	s := block.Find(".baav")
	if s.Size() == 0 {
		// Chinese => English
		pronunciation = append(pronunciation, block.Find(".wordbook-js .phonetic").Text())
	} else {
		// English => Chinese
		s.Children().Each(func(i int, s *goquery.Selection) {
			// s.Text() will return both self's text and childrens text
			// so here uses s.Contents().Get(0).Data to get text
			pronounce := strings.TrimSpace(s.Contents().Get(0).Data)
			phonetic := s.Find(".phonetic").Text()
			pronunciation = append(pronunciation, fmt.Sprintf("%s %s", pronounce, phonetic))
		})
	}
	return pronunciation
}

func (dict youdao) parseExplanation(block *goquery.Selection) []map[string]string {
	explanation := make([]map[string]string, 0, 5)
	children := block.Find(".trans-container > ul").Children()
	if children.First().Is("li") {
		children.Each(func(i int, s *goquery.Selection) {
			data := strings.SplitAfterN(s.Text(), ".", 2)
			explain := map[string]string{
				"prop":    data[0],
				"meaning": strings.TrimSpace(data[1]),
			}
			explanation = append(explanation, explain)
		})
	} else {
		if _, exists := children.First().Children().First().Attr("style"); exists {
			children.Each(func(i int, s *goquery.Selection) {
				prop := s.Children().First().Text()
				meanings := make([]string, 0, 10)
				s.Find(".search-js").Each(func(i int, s *goquery.Selection) {
					meanings = append(meanings, s.Text())
				})
				explain := map[string]string{
					"prop":    prop,
					"meaning": strings.Join(meanings, "; "),
				}
				explanation = append(explanation, explain)
			})
		} else {
			children.Each(func(i int, s *goquery.Selection) {
				explain := map[string]string{
					"prop":    "",
					"meaning": s.Find(".search-js").Text(),
				}
				explanation = append(explanation, explain)
			})
		}
	}
	return explanation
}

func (dict youdao) getIndention(explanation []map[string]string) int {
	length := 0
	for _, explain := range explanation {
		tmpL := utf8.RuneCountInString(explain["prop"])
		if tmpL > length {
			length = tmpL
		}
	}
	return length
}

func (dict youdao) display(word string, pronunciation []string, explanation []map[string]string) {
	fmt.Println(word)
	for _, p := range pronunciation {
		fmt.Printf("%s  ", p)
	}
	fmt.Print("\n")
	indention := dict.getIndention(explanation)
	for _, explain := range explanation {
		if indention != 0 {
			fmt.Printf("%[3]*[1]s %[2]s\n", explain["prop"], explain["meaning"], indention)
		} else {
			fmt.Println(explain["meaning"])
		}
	}
}
