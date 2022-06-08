package main

import (
	"flag"
	"fmt"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"path/filepath"
)

func ReadAll(path string) ([][]byte, error) {
	res := new([][]byte)
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, d := range dir {
		if !d.IsDir() && filepath.Ext(d.Name()) == ".json" {
			file, err := ioutil.ReadFile(path + "/" + d.Name())
			if err == nil {
				*res = append(*res, file)
			}
		}
	}
	return *res, nil
}

func main() {
	var pathToJson string
	flag.StringVar(&pathToJson, "p", "./resources", "path to files json")
	flag.Parse()

	jsons, err := ReadAll(pathToJson)
	if err != nil {
		fmt.Println(err)
		return
	}

	sc, err := stan.Connect("test-cluster", "master")
	if err != nil {
		fmt.Println(err)
	}
	defer func(sc stan.Conn) {
		err := sc.Close()
		if err != nil {

		}
	}(sc)

	// Simple Synchronous Publisher
	for i, m := range jsons {
		err = sc.Publish("test", m) // does not return until an ack has been received from NATS Streaming
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Published ok:", i)
		}
	}
}
