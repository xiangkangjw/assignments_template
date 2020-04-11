package cos418_hw1_1

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Find the top K most common words in a text document.
// 	path: location of the document
//	numWords: number of words to return (i.e. k)
//	charThreshold: character threshold for whether a token qualifies as a word,
//		e.g. charThreshold = 5 means "apple" is a word but "pear" is not.
// Matching is case insensitive, e.g. "Orange" and "orange" is considered the same word.
// A word comprises alphanumeric characters only. All punctuations and other characters
// are removed, e.g. "don't" becomes "dont".
// You should use `checkError` to handle potential errors.
func topWords(path string, numWords int, charThreshold int) []WordCount {
	// TODO: implement me
	// HINT: You may find the `strings.Fields` and `strings.ToLower` functions helpful
	// HINT: To keep only alphanumeric characters, use the regex "[^0-9a-zA-Z]+"
	file, err := os.Open(path)
	checkError(err)
	reg, err := regexp.Compile("[^0-9a-zA-Z]+")
	checkError(err)
	words := readFile(file, reg, charThreshold)
	wordMap := make(map[string]int)
	for _, word := range words {
		if count, ok := wordMap[word]; ok {
			wordMap[word] = count + 1
		} else {
			wordMap[word] = 1
		}
	}
	wordCounts := make([]WordCount, len(wordMap))
	i := 0
	for k, v := range wordMap {
		wordCounts[i] = WordCount{k, v}
		i++
	}
	sortWordCounts(wordCounts)
	return wordCounts[:numWords]
}

func readFile(file *os.File, reg *regexp.Regexp, charThreshold int) []string {
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	result := make([]string, 0)
	for scanner.Scan() {
		word := scanner.Text()
		word = strings.ToLower(word)
		word = reg.ReplaceAllString(word, "")
		if len(word) < charThreshold {
			continue
		}
		result = append(result, word)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.
// DO NOT MODIFY THIS FUNCTION!
func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}
