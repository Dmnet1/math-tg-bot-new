package tgbot

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func UnmarshalEquations() ([]Equation, error) { // поправил название
	file, err := os.Open("types.json")
	if err != nil {
		log.Panic(err)
	} else {
		log.Println(file)
	}

	defer file.Close() // поставил сразу после открытия

	/*fi, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}*/

	fi, err := io.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}

	/*var data = make([]byte, fi.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Panic(err)
	}*/

	var result []Equation
	err = json.Unmarshal(fi, &result)

	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}
