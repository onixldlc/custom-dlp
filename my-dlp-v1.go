package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
)

func main() {
    // Define supported modes
    modes := map[string][]string{
        "mp3":     {"-x", "--audio-format", "mp3", "--audio-quality", "0"},
        "mp4":     {"-f", "bestvideo+bestaudio/best[ext=mp4]"},
        "no-audio": {"-f", "bestvideo"},
        "musvid":  {"-f", "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]"},
    }

    // Parse mode and URL from command-line arguments
    mode := flag.String("mode", "", "Mode to run dlp: mp3, mp4, no-audio, musvid")
    url := flag.String("url", "", "URL to the video")
    flag.Parse()

    // Validate mode and URL input
    if *mode == "" || *url == "" {
        fmt.Println("Usage:")
        fmt.Println("  -mode string")
        fmt.Println("        Mode to run dlp: mp3, mp4, no-audio, musvid")
        fmt.Println("  -url string")
        fmt.Println("        URL to the video")
        fmt.Println("\nExample:")
        fmt.Println("  go run main.go -mode=mp4 -url=https://example.com/video")
        os.Exit(1)
    }

    // Check if mode is supported
    options, exists := modes[*mode]
    if !exists {
        fmt.Printf("Mode '%s' not recognized.\n", *mode)
        fmt.Println("Available modes: mp3, mp4, no-audio, musvid")
        os.Exit(1)
    }

    // Construct yt-dlp command
    args := append(options, *url)
    cmd := exec.Command("yt-dlp", args...)

    // Run the command
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    fmt.Printf("Running yt-dlp in %s mode for URL: %s\n", *mode, *url)
    if err := cmd.Run(); err != nil {
        fmt.Printf("Error executing yt-dlp: %v\n", err)
        os.Exit(1)
    }
}

