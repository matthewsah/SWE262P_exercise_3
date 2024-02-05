package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

// TFTheOne struct
type TFTheOne struct {
	value interface{}
}

// bind the current value to a function and return TFTheOne
func (tfo *TFTheOne) Bind(f func(interface{}) interface{}) *TFTheOne {
	tfo.value = f(tfo.value)
	return tfo
}

// print the value
func (tfo *TFTheOne) PrintMe() {
	fmt.Println(tfo.value)
}

func ReadFile(path interface{}) interface{} {
	filePath, ok := path.(string)
	if !ok {
		panic("Invalid input")
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func FilterChars(strData interface{}) interface{} {
	pattern := regexp.MustCompile(`[\W_]+`)
	return pattern.ReplaceAllString(strData.(string), " ")
}

func Normalize(strData interface{}) interface{} {
	return strings.ToLower(strData.(string))
}

func Scan(strData interface{}) interface{} {
	return strings.Split(strData.(string), " ")
}

func RemoveStopWords(wordList interface{}) interface{} {
	stopWordsData, err := ioutil.ReadFile("../stop_words.txt")
	if err != nil {
		panic(err)
	}
	stopWords := strings.Split(strings.TrimSpace(string(stopWordsData)), ",")
	words := wordList.([]string)
	result := []string{}
	for _, w := range words {
		found := false
		for _, sw := range stopWords {
			if w == sw {
				found = true
				break
			}
		}
		if !found && len(w) > 1 {
			result = append(result, w)
		}
	}
	return result
}

func Frequencies(wordList interface{}) interface{} {
	words := wordList.([]string)
	freqs := make(map[string]int)
	for _, word := range words {
		freqs[word]++
	}
	return freqs
}

func Sort(wordFreqs interface{}) interface{} {
	freqs := wordFreqs.(map[string]int)
	sortedFreqs := make([][2]interface{}, 0, len(freqs))
	for k, v := range freqs {
		sortedFreqs = append(sortedFreqs, [2]interface{}{k, v})
	}
	sort.Slice(sortedFreqs, func(i, j int) bool {
		return sortedFreqs[i][1].(int) > sortedFreqs[j][1].(int)
	})
	return sortedFreqs
}

func Top25Freqs(wordFreqs interface{}) interface{} {
	freqs := wordFreqs.([][2]interface{})
	var result strings.Builder
	for _, tf := range freqs[:25] {
		result.WriteString(fmt.Sprintf("%s - %d\n", tf[0], tf[1]))
	}
	return result.String()
}

func main() {
	T := TFTheOne{value: os.Args[1]}
	T.
		Bind(ReadFile).
		Bind(FilterChars).
		Bind(Normalize).
		Bind(Scan).
		Bind(RemoveStopWords).
		Bind(Frequencies).
		Bind(Sort).
		Bind(Top25Freqs).
		PrintMe()
}
