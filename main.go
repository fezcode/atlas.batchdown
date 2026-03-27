package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var Version = "dev"

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	arg := os.Args[1]
	if arg == "-v" || arg == "--version" {
		fmt.Printf("atlas.batchdown v%s\n", Version)
		return
	}
	if arg == "-h" || arg == "--help" || arg == "help" {
		showHelp()
		return
	}

	filePath := arg
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found.\n", filePath)
		os.Exit(1)
	}

	urls, err := readLines(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting batch downloads for %d URLs...\n", len(urls))
	for i, url := range urls {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		fmt.Printf("[%d/%d] Downloading %s...\n", i+1, len(urls), url)
		err := downloadFile(url)
		if err != nil {
			fmt.Printf("  Error downloading %s: %v\n", url, err)
		} else {
			fmt.Printf("  Successfully downloaded.\n")
		}
	}
	fmt.Println("All tasks completed.")
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func downloadFile(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	filename := getFilename(url, resp)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getFilename(url string, resp *http.Response) string {
	// Try to get filename from URL
	base := filepath.Base(url)
	if base == "." || base == "/" {
		// Fallback to timestamp if URL doesn't have a clear filename
		return fmt.Sprintf("download_%d", time.Now().Unix())
	}
	
	// Remove query params if any
	if idx := strings.Index(base, "?"); idx != -1 {
		base = base[:idx]
	}
	
	return base
}

func showHelp() {
	fmt.Println("Atlas BatchDown - A powerful batch downloader for the Atlas Suite.")
	fmt.Println("\nUsage:")
	fmt.Println("  atlas.batchdown <file.txt>     Download all links from the text file")
	fmt.Println("  atlas.batchdown -h, --help     Show this help information")
	fmt.Println("  atlas.batchdown -v, --version  Show version info")
	fmt.Println("\nFile Format:")
	fmt.Println("  The input file should contain one URL per line.")
}
