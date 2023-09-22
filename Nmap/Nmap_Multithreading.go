package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("")
		return
	}

	filename := os.Args[1]
	threads, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("Threads Inválidos:", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error al abrir el archivo:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, threads)

	progress := uiprogress.New()
	progress.Start()
	bar := progress.AddBar(100)
	bar.AppendCompleted()
	bar.PrependElapsed()

	for scanner.Scan() {
		host := scanner.Text()

		wg.Add(1)
		semaphore <- struct{}{} 
		go func(host string) {
			defer func() {
				<-semaphore 
				wg.Done()
				bar.Incr()
			}()

			outputFile := host + ".txt" 
			cmd := exec.Command("nmap", "-sCV", "--min-rate", "10000", "-Pn", "-oN", outputFile, host)
			err := cmd.Run()
			if err != nil {
				log.Printf("Error en el escaneo de Nmap de %s: %s\n", host, err)
				return
			}

			time.Sleep(500 * time.Millisecond)
		}(host)
	}

	wg.Wait()
	progress.Stop()

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file:", err)
	}
}
