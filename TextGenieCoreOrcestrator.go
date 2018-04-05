package ss

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/orcaman/concurrent-map"
)

type match struct {
	inputItem     string
	standardItems []string
	computePairs  [][]string
	Response      cmap.ConcurrentMap
}

func (mp match) createPairs() [][]string {
	inputCategory := strings.Join(strings.Fields(mp.inputItem), " ")
	computePairs := [][]string{}
	for _, standardCategoryItems := range mp.standardItems {
		standardCategoryItemsCleaned := strings.Join(strings.Fields(standardCategoryItems), " ")
		computePairs = append(computePairs, []string{inputCategory, standardCategoryItemsCleaned})

	}
	return computePairs

}

func (mp match) threadChild(computePair []string, wg *sync.WaitGroup) {
	defer wg.Done()
	ip := ItemPair{}
	score := (ip.TextGenieProcessor(computePair[0], computePair[1]))
	mp.Response.Set(computePair[1], score)
}

func (mp match) scoreByCategory() string {
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

//MatchProcessor
func (mp match) MatchProcessor(inputItem string, standardItems []string) string {
	mp.inputItem = inputItem
	mp.standardItems = standardItems
	mp.computePairs = mp.createPairs()
	return mp.scoreByCategory()

}
func main() {

	mp := match{}

	fmt.Println(mp.MatchProcessor("sta", []string{"standard", "elephant", "andard"}))

}
