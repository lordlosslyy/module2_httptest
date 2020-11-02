package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := flag.Int("port", 8080, "Port to accept connection on")
	host := flag.String("host", "127.0.0.1", "Host or IP to bind to")
	requestsFilename := flag.String("requestFile", "requests.txt", "File for requests")
	goldensFilename := flag.String("goldenFile", "goldenheaders.txt", "File for requests")
	delim := flag.String("delimeter", "\r\n", "Delimeter that seperates golden")
	flag.Parse()

	// check(filename)
	println(*host + ":" + strconv.Itoa(*port))

	requests := readrequestsFile(*requestsFilename, *delim) // read requests
	goldens := readgoldensFile(*goldensFilename)            // read golden answer

	fmt.Println(goldens)
	// fmt.Println(reflect.TypeOf(goldens))
	// fmt.Println(reflect.TypeOf(requests))

	for i := 0; i < len(requests); i++ {
		// Connect to the server
		conn, err := net.Dial("tcp", *host+":"+strconv.Itoa(*port))
		if err != nil {
			log.Panicln(err)
		}
		defer conn.Close()

		// Sent the request to the server
		sentbytes := []byte(requests[i])
		_, err = conn.Write(sentbytes)

		if err != nil {
			log.Panic(err)
			continue
		}

		// Read the response from the server
		remaining := ""
		count := 0
		anwser := ""
		//rest := ""
		//context := ""
		for {
			buf := make([]byte, 10)
			size, err := conn.Read(buf)
			if err != nil {
				log.Println("err != nil")
				break
			}
			// log.Println(size)
			data := buf[:size]
			remaining = remaining + string(data)
			//rest = remaining
			for strings.Contains(remaining, "\r\n") {

				idx := strings.Index(remaining, "\r\n")
				//context = remaining[:idx]
				//log.Println(remaining[:idx])

				if count == 0 {
					anwser = remaining[:idx]
					//log.Println(anwser)
				}
				remaining = remaining[idx+len("\r\n"):]
				count++
			}
		}
		// log.Println(rest)
		/*
			if i == 0 {
				fmt.Println(context)
			}
		*/
		// check context
		/*
			if i == 0 && check(context, ) {
				fmt.Printf("First", i+1)
			}
		*/
		// Check the response equal to the golden answer
		if check(anwser, goldens[i]) {
			fmt.Printf("%v. Pass\n", i+1)
		} else {
			fmt.Printf("%v. Wrong\n", i+1)
		}
	}
}

// Read Requests File
func readrequestsFile(filename string, delim string) []string {
	log.Println("Reading from file " + filename)

	f, err := os.Open(filename) // read file
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Panicln(err)
		}
	}()

	remaining := ""
	requests := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		buf := make([]byte, 10)
		buf = s.Bytes()
		remaining = remaining + string(buf) + delim

		for strings.Contains(remaining, "$") {
			idx := strings.Index(remaining, "$")
			log.Println(remaining[:idx])
			requests = append(requests, remaining[:idx])
			remaining = remaining[idx+len("$"):]
		}
	}

	err = s.Err()
	if err != nil {
		log.Panicln(err)
	}

	return requests
}

// Read Golden File
func readgoldensFile(filename string) []string {
	log.Println("Reading from file " + filename)

	f, err := os.Open(filename) // read file
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Panicln(err)
		}
	}()

	goldens := []string{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		goldens = append(goldens, s.Text())
	}

	err = s.Err()
	if err != nil {
		log.Panicln(err)
	}

	return goldens
}

func check(response, golden string) bool {

	return response == golden
}
