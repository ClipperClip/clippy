package cli

import (
	"fmt"
	"github.com/evanweissburg/clippy/pkg/client"
	"github.com/mholt/archiver/v3"
	"io"
	"log"
	"os"
)

func invalid_usage() {
	fmt.Println("Correct usage: <>")
	os.Exit(1)
}

func Execute() {
	if len(os.Args) != 3 {
		invalid_usage()
	}

	switch os.Args[1] {
	case "put":
		filename := os.Args[2]
		put(filename)

	case "get":
		clipcode := os.Args[2]
		get(clipcode)

	default:
		invalid_usage()
	}
}

func put(filename string) {
	err := archiver.Archive([]string{filename}, ".clip.zip")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(".clip.zip")
	if err != nil {
		log.Fatal(err)
	}

	cl := client.Client{
		BaseURL: "http://localhost:8080/",
	}

	clipcode, err := cl.Upload(file)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Recieved clipcode %s\n", clipcode)

	err = os.Remove(".clip.zip")
	if err != nil {
		log.Fatalln(err)
	}
}

func get(clipcode string) {
	cl := client.Client{
		BaseURL: "http://localhost:8080/",
	}

	data, err := cl.Download(clipcode)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Create(".clip.zip")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(file, data)
	err = file.Close()
	if err != nil {
		log.Fatalln(err)
	}

	err = archiver.Unarchive(".clip.zip", ".")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Remove(".clip.zip")
	if err != nil {
		log.Fatalln(err)
	}
}
