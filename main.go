package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var BASE_URL string = "https://pokeapi.co/api/v2/pokemon/"

func NewSlice(start, end, step int) []int {
	if step <= 0 || end < start {
		return []int{}
	}
	s := make([]int, 0, 1+(end-start)/step)
	for start <= end {
		s = append(s, start)
		start += step
	}
	return s
}

func getPokemon(index int, outChannel chan<- []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, BASE_URL+strconv.Itoa(index+1), nil)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%v\n", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}
	// fmt.Printf("%v\n", string(body))
	outChannel <- body
}

type Pokemon struct {
	Name string `json:"name"`
}

func getPokemons() {
	type Result struct {
		Size int
		Time int64
	}

	result := make([]Result, 0)

	for _, i := range NewSlice(10, 800, 20) {
		fmt.Printf("Getting %v pokemons\n", i)
		start := time.Now()
		wg := sync.WaitGroup{}
		outChannel := make(chan []byte, i)
		for j := 0; j < i; j++ {
			wg.Add(1)
			go getPokemon(j, outChannel, &wg)
		}
		wg.Wait()
		close(outChannel)

		elapsed := time.Since(start)
		result = append(result, Result{i, elapsed.Milliseconds()})

		log.Printf("Getting pokemons took %s", elapsed)

		result := make([]string, 0, i)
		for r := range outChannel {
			var p Pokemon
			jsonErr := json.Unmarshal(r, &p)
			// fmt.Printf("%v\n", string(r))
			if jsonErr != nil {
				err := fmt.Errorf("%v", string(r))
				panic(err)
			}
			result = append(result, p.Name)
		}
		for _, name := range result[:10] {
			fmt.Printf("%v ", name)
		}
		fmt.Print("\n")
	}

	var data [][]string
	for _, r := range result {
		data = append(data, []string{strconv.Itoa(r.Size), strconv.FormatInt(r.Time, 10)})
	}
	file, err := os.Create("golang_async.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	w.WriteAll(data)
	w.Flush()
}

func main() {
	getPokemons()
}
