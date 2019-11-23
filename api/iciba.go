package api

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

type iciba struct{}

const icibaURL = "https://www.iciba.com/%s"

// Search looks up "word" from www.iciba.com and prints the result on screen
func (dict iciba) Search(word string) {
	block, err := dict.getContent(word)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	pronunciation := dict.parsePronunciation(block)
	explanation := dict.parseExplanation(block)
	dict.display(word, pronunciation, explanation)
}

func (dict iciba) getContent(word string) (*goquery.Selection, error) {
	url := fmt.Sprintf(icibaURL, word)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	block := doc.Find(".in-base")
	return block, nil
}

func (dict iciba) parsePronunciation(block *goquery.Selection) []string {
	pronunciation := make([]string, 0, 2)
	block.Find(".base-speak > span").Each(func(i int, s *goquery.Selection) {
		pronunciation = append(pronunciation, s.Find("span").Text())
	})
	return pronunciation
}

func (dict iciba) parseExplanation(block *goquery.Selection) map[string]string {
	explanation := make(map[string]string)
	block.Find(".base-list > .clearfix").Each(func(i int, s *goquery.Selection) {
		var meaning strings.Builder
		prop := s.Find(".prop").Text()
		s.Find("p > span").Each(func(i int, s *goquery.Selection) {
			meaning.WriteString(s.Text())
		})
		explanation[prop] = meaning.String()
	})
	return explanation
}

func (dict iciba) getIndention(explanation map[string]string) int {
	length := 0
	for prop := range explanation {
		tmpL := utf8.RuneCountInString(prop) // len() does not work as expected for CJK charactors
		if tmpL > length {
			length = tmpL
		}
	}
	return length
}

func (dict iciba) display(word string, pronunciation []string, explanation map[string]string) {
	fmt.Println(word)
	for _, p := range pronunciation {
		fmt.Printf("%s  ", p)
	}
	fmt.Print("\n")
	indention := dict.getIndention(explanation)
	for k, v := range explanation {
		fmt.Printf("%[3]*[1]s %[2]s\n", k, v, indention)
	}
}
