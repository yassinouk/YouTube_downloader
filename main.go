package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func runCommand(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Set the working directory to /root/downloads
	if err := os.Chdir("/media/videobucket"); err != nil {
		fmt.Printf("Error changing working directory: %v\n", err)
		return
	}

	command := fmt.Sprintf("yt-dlp %s", url)
	cmd := exec.Command("sh", "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing command for URL %s: %v\n", url, err)
	}
}

func main() {
	// Open the CSV file
	file, err := os.Open("urls.csv")
	if err != nil {
		fmt.Printf("Error opening CSV file: %v\n", err)
		return
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Slice to store the URLs
	var urls []string

	// Read and process each row in the CSV file
	for {
		// Read the next row from the CSV
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}

		// Assuming the URL is in the first column (index 0)
		url := record[0]

		// Append the URL to the slice
		urls = append(urls, url)
	}

	// Concurrently execute commands for each URL
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go runCommand(url, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
