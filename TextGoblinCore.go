package main

import (
	"strings"
	"github.com/jaideepkekre/goabber"
)

import "github.com/charlesvdv/fuzmatch"
import "github.com/antzucaro/matchr"

type ItemPair struct {
	inputItems    string
	standardItems string
	Score         int
}

func (ip ItemPair) abberBoost(inputItem string, standardItem string) int {
	if goabber.Abber(inputItem, standardItem) {
		return 602
	}
	return 0
}

func metaphoneScorer(inputMetaPhone string, stdMetaPhone string) int {

	metaLevenstienDistance := matchr.Levenshtein(inputMetaPhone, stdMetaPhone)

	if metaLevenstienDistance <= len(stdMetaPhone)/2 {

		score := fuzmatch.Ratio(inputMetaPhone, stdMetaPhone)
		return score
	} else {
		return 0

	}

}

//FuzzPhoneticScorer returns the best score scaled out of 1000 using fuzzy, phonetic algorithms
func (ip ItemPair) fuzzPhoneticCalulator(inputItem string, standardItem string) int {
	if inputItem[0] != standardItem[0] {
		return 0
	}
	fuzzPartialRatio := fuzmatch.PartialRatio(inputItem, standardItem)
	fuzzLevenstienDistance := matchr.Levenshtein(inputItem, standardItem)
	inputMetaPhone1, inputMetaPhone2 := matchr.DoubleMetaphone(inputItem)
	stdMetaPhone1, stdMetaPhone2 := matchr.DoubleMetaphone(standardItem)

	metaPhoneScore1 := metaphoneScorer(inputMetaPhone1, stdMetaPhone1)
	metaPhoneScore2 := metaphoneScorer(inputMetaPhone1, stdMetaPhone2)
	metaPhoneScore3 := metaphoneScorer(inputMetaPhone2, stdMetaPhone1)
	metaPhoneScore4 := metaphoneScorer(inputMetaPhone2, stdMetaPhone2)

	fuzzPartialRatio = fuzzPartialRatio * 10
	metaPhoneScore1 = metaPhoneScore1 * 10
	metaPhoneScore2 = metaPhoneScore2 * 10
	metaPhoneScore3 = metaPhoneScore3 * 10
	metaPhoneScore4 = metaPhoneScore4 * 10
	abberScore := ip.abberBoost(inputItem, standardItem)

	if fuzzLevenstienDistance > len(standardItem) {
		fuzzPartialRatio = 0
	}

	scores := []int{}
	scores = append(scores, fuzzPartialRatio, metaPhoneScore1,
		metaPhoneScore2, metaPhoneScore3,
		metaPhoneScore4)
	max := -1
	for _, miniScore := range scores {
		if max < miniScore {
			max = miniScore
		}
	}
	max = max + abberScore
	if max > 600 {
		return max
	} else {
		return 0
	}

}



func (ip ItemPair) TextGenieProcessor(inputItems string, standardItems string) int {

	ip.inputItems = inputItems
	ip.standardItems = standardItems

	for _, intputItemsSlice := range strings.Split(ip.inputItems, " ") {
		for _, standardItemsSlice := range strings.Split(ip.standardItems, " ") {
			ip.Score = ip.Score + ip.fuzzPhoneticCalulator(intputItemsSlice, standardItemsSlice)
		}

	}
	// ip.fuzzPhoneticCalulator()
	return ip.Score

}
