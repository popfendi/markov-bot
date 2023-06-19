package markov

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mb-14/gomarkov"
)

var MC *gomarkov.Chain

func Init() {

	chain, err := loadModel()
	if err != nil {
		chain = buildModel(1)
	}

	MC = chain

	go saveChainToFileAtInterval()

}

func saveChainToFileAtInterval() {
	for {
		log.Println("Saving Chain to 'model.json'")
		//The chain is JSON serializable
		jsonObj, _ := json.Marshal(MC)
		jsonErr := ioutil.WriteFile("model.json", jsonObj, 0644)
		if jsonErr != nil {
			fmt.Println(jsonErr)
		}
		time.Sleep(300 * time.Second)
	}
}

func buildModel(order int) *gomarkov.Chain {
	chain := gomarkov.NewChain(order)
	for _, data := range getDataset("init2.txt") {
		chain.Add(Split(data))
	}
	return chain
}

func loadModel() (*gomarkov.Chain, error) {
	var chain gomarkov.Chain
	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		return &chain, err
	}
	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}
	return &chain, nil
}

func getDataset(fileName string) []string {
	file, _ := os.Open(fileName)
	scanner := bufio.NewScanner(file)
	var list []string
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	return list
}

func Split(str string) []string {
	return strings.Split(str, " ")
}

func Generate() string {
	order := MC.Order
	tokens := make([]string, 0)
	for i := 0; i < order; i++ {
		tokens = append(tokens, gomarkov.StartToken)
	}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := MC.Generate(tokens[(len(tokens) - order):])
		tokens = append(tokens, next)
	}

	s := strings.Join(tokens[order:len(tokens)-1], " ")

	if len(strings.TrimSpace(s)) == 0 {
		s = Generate()
	}

	return s
}

func Train(text string) {
	MC.Add(Split(text))
}
