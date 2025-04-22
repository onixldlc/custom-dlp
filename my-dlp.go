package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func main() {
    // Define supported modes and their corresponding yt-dlp options
    modes := map[string][]string{
        "mp3":      {"-x", "--audio-format", "mp3", "--audio-quality", "0"},
        "mp4":      {"-f", "bestvideo+bestaudio/best[ext=mp4]", "-S", "vcodec:h264,res,acodec:m4a"},
        "no-audio": {"-f", "bestvideo/best[ext=mp4]", "-S", "vcodec:h264,res,acodec:m4a"},
        "musvid":   {"-f", "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]", "-S", "vcodec:h264,res,acodec:m4a"},
    }

    // Define flags
    modeFlag := flag.String("mode", "", "Mode to run dlp: mp3, mp4, no-audio, musvid")
    urlFlag := flag.String("url", "", "URL to the video")

    // Customize the usage message
    flag.Usage = func() {
        usageText := `
my-dlp 1.0.3
Usage:
  [mode] -url <URL>
  -mode=<mode> -url=<URL>

Modes:
  mp3       Only grab the mp3 
            Example: -mode=mp3 -url=https://example.com/video
  mp4       Grab video + audio output always mp4
            Example: -mode=mp4 -url=https://example.com/video
  no-audio  Video only
            Example: -mode=no-audio -url=https://example.com/video
  musvid    Limit the video to best audio
            Example: -mode=musvid -url=https://example.com/video

Positional Arguments:
  mode       (optional) Mode to run dlp: mp3, mp4, no-audio, musvid
  url        URL to the video

Examples:
  my-dlp -mode=mp3 -url=https://example.com/audio
  my-dlp mp4 https://example.com/video
  my-dlp https://example.com/video

`
        fmt.Fprintf(flag.CommandLine.Output(), "%s\n", usageText)
    }

    // Parse flags
    flag.Parse()

    var mode string
    var url string

    args := flag.Args()

    // Determine mode and URL based on flags and positional arguments
    if *modeFlag != "" {
        mode = *modeFlag
        if *urlFlag != "" {
            url = *urlFlag
        } else if len(args) > 0 {
            url = args[0]
        }
    } else {
        if len(args) > 0 && isValidMode(args[0], modes) {
            mode = args[0]
            if len(args) > 1 {
                url = args[1]
            }
        } else {
            if *urlFlag != "" {
                url = *urlFlag
            } else if len(args) > 0 {
                url = args[0]
            }
        }
    }

    // Validate URL
    if url == "" {
        fmt.Println()
        fmt.Println("Error: URL to the video is required.")
        fmt.Println()
        flag.Usage()
        os.Exit(1)
    }

    // Validate mode if provided
    var options []string
    if mode != "" {
        var exists bool
        options, exists = modes[strings.ToLower(mode)]
        if !exists {
            fmt.Printf("Mode '%s' not recognized.\n", mode)
            fmt.Println("Available modes: mp3, mp4, no-audio, musvid")
            os.Exit(1)
        }
    }

    // Construct yt-dlp command
    argsForYtDlp := []string{}
    if len(options) > 0 {
        argsForYtDlp = append(argsForYtDlp, options...)
    }
    argsForYtDlp = append(argsForYtDlp, url)

    // Display the mode and URL being used
    if mode != "" {
        fmt.Printf("Running yt-dlp in '%s' mode for URL: %s\n", mode, url)
    } else {
        fmt.Printf("Running yt-dlp in default mode for URL: %s\n", url)
    }

    // Execute yt-dlp
    cmd := exec.Command("yt-dlp", argsForYtDlp...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        fmt.Printf("Error executing yt-dlp: %v\n", err)
        os.Exit(1)
    }
}

// isValidMode checks if the provided mode is supported
func isValidMode(mode string, modes map[string][]string) bool {
    _, exists := modes[strings.ToLower(mode)]
    return exists
}

