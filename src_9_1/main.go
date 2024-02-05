package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

type pair struct {
	word  string
	count int
}

func readFile(path string, nextFunc func([]rune, func([]string, func([]string, func([]string, func([]pair, func([]pair))))))) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	charData := make([]rune, 0)

	for scanner.Scan() {
		line := scanner.Text() + " "
		for _, char := range line {
			charData = append(charData, char)
		}
	}

	nextFunc(charData, scan)
}

func filterCharsAndNormalize(charData []rune, nextFunc func([]string, func([]string, func([]string, func([]pair, func([]pair)))))) {
	stringData := make([]string, 0)

	for idx, char := range charData {
		if !(rune('a') <= char && char <= rune('z') || rune('A') <= char && char <= rune('Z') || rune('0') <= char && char <= rune('9')) {
			charData[idx] = ' '
		} else {
			charData[idx] = unicode.ToLower(char)
		}
	}

	for _, char := range charData {
		stringData = append(stringData, string(char))
	}

	nextFunc(stringData, removeStopwords)
}

func scan(stringData []string, nextFunc func([]string, func([]string, func([]pair, func([]pair))))) {
	text := strings.Join(stringData, "")
	words := strings.Split(text, " ")
	nextFunc(words, frequencies)
}

func removeStopwords(words []string, nextFunc func([]string, func([]pair, func([]pair)))) {
	stopwords, err := os.ReadFile("../stop_words.txt")
	if err != nil {
		log.Fatal(err)
	}
	stopwordsList := strings.Split(string(stopwords), ",")
	asciiVal := 97
	for i := 0; i < 26; i++ {
		char := string(rune(asciiVal + i))
		stopwordsList = append(stopwordsList, char)
	}
	stopwordsList = append(stopwordsList, string(""))

	stopwordIdxs := make([]int, 0)
	for idx, word := range words {
		for _, stopword := range stopwordsList {
			if word == stopword {
				stopwordIdxs = append(stopwordIdxs, idx)
				break
			}
		}
	}

	for i := len(stopwordIdxs) - 1; i >= 0; i-- {
		idx := stopwordIdxs[i]
		left := words[:idx]
		right := words[idx+1:]
		words = append(left, right...)
	}

	nextFunc(words, sortFreqs)
}

func frequencies(words []string, nextFunc func([]pair, func([]pair))) {
	wordFreqs := make([]pair, 0)

	for _, word := range words {
		found := false
		for idx, pair := range wordFreqs {
			if word == pair.word {
				wordFreqs[idx].count += 1
				found = true
				break
			}
		}
		if !found {
			wordFreqs = append(wordFreqs, pair{word, 1})
		}
	}

	nextFunc(wordFreqs, printFreqs)
}

func sortFreqs(wordFreqs []pair, nextFunc func([]pair)) {
	sort.SliceStable(wordFreqs, func(p1, p2 int) bool {
		return wordFreqs[p1].count > wordFreqs[p2].count
	})
	nextFunc(wordFreqs)
}

func printFreqs(wordFreqs []pair) {
	for _, pair := range wordFreqs[:25] {
		fmt.Println(pair.word, " - ", pair.count)
	}
}

func main() {
	readFile(os.Args[1], filterCharsAndNormalize)
}
