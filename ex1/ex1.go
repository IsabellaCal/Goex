package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type exercize struct {
	question string
	solution string
}

var score int

// crea l'esercizio
func createExercize(data [][]string) []exercize {
	var questions []exercize
	for _, line := range data {
		questions = append(questions, exercize{question: line[0], solution: line[1]})
	}
	return questions
}

// chiede domande
func ask(exercizes []exercize) {
	for _, exercize := range exercizes {
		fmt.Println(exercize.question)
		var answer string
		fmt.Scanln(&answer)
		if answer == exercize.solution {
			score++
		}
	}
}

func main() {
	var name string
	var timer int
	flag.StringVar(&name, name, "problems.csv", "nome del file")
	flag.IntVar(&timer, "tempo", 30, "timer")
	flag.Parse()
	//apro il file

	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	// leggo i dati
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	exercizes := createExercize(data)

	//timer
	done := make(chan bool)
	go func() {
		ask(exercizes)
		done <- true
	}()
	select {
	case <-done:
		fmt.Println("Il tuo punteggio è: ", score, "/", len(exercizes))
	case <-time.After(time.Duration(timer) * time.Second):
		println("Tempo scaduto!")
		fmt.Println("Il tuo punteggio è: ", score, "/", len(exercizes))
	}

}
