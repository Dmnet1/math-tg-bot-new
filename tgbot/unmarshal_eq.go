package tgbot

import (
	"encoding/xml"
	"log"
	"os"
)

func UnmarshalEquations() (Equations, error) { // поправил название
	file, err := os.Open("types.xml")
	if err != nil {
		log.Panic(err)
	} else {
		log.Println(file)
	}

	defer file.Close() // поставил сразу после открытия

	fi, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}

	var data = make([]byte, fi.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Panic(err)
	}

	var v Equations
	err = xml.Unmarshal(data, &v)

	if err != nil {
		log.Println(err)
		return v, err
	}
	return v, nil
}
