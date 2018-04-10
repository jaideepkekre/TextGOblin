package textgoblin

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/orcaman/concurrent-map"
)

type Match struct {
	inputItem     string
	standardItems []string
	computePairs  [][]string
	Response      cmap.ConcurrentMap
}

func (mp Match) createPairs() [][]string {
	inputCategory := strings.Join(strings.Fields(mp.inputItem), " ")
	computePairs := [][]string{}
	for _, standardCategoryItems := range mp.standardItems {
		standardCategoryItemsCleaned := strings.Join(strings.Fields(standardCategoryItems), " ")
		computePairs = append(computePairs, []string{inputCategory, standardCategoryItemsCleaned})

	}
	return computePairs

}

func (mp Match) threadChild(computePair []string, wg *sync.WaitGroup) {
	defer wg.Done()
	ip := ItemPair{}
	score := (ip.TextGenieProcessor(computePair[0], computePair[1]))
	mp.Response.Set(computePair[1], score)
}

func (mp Match) scoreByCategory() string {
	mp.Response = cmap.New()
	var wg sync.WaitGroup
	for _, computePair := range mp.computePairs {
		wg.Add(1)
		go mp.threadChild(computePair, &wg)
	}
	wg.Wait()
	jsonResp, _ := json.Marshal(mp.Response)
	return string(jsonResp)

}

//
func (mp Match) MatchProcessor(inputItem string, standardItems []string) string {
	mp.inputItem = inputItem
	mp.standardItems = standardItems
	mp.computePairs = mp.createPairs()
	return mp.scoreByCategory()

}


