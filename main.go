package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	arg "github.com/alexflint/go-arg"
)

const (
	GOARCH string = runtime.GOARCH
	GOOS   string = runtime.GOOS
)

var args struct {
	User         string `arg:"required"`
	Project      string `arg:"required"`
	Kernel       string `arg:"" default:"" help:"linux/darwin/window (Default is your running Kernel)"`
	Architecture string `arg:"" default:"" help:"auto/amd64/386"`
	Lastest      bool   `arg:"" help:"Get last release (Default is Latest Stable Release)"`
	Download     bool   `arg:"" help:"Download file"`
}

func main() {
	arg.MustParse(&args)

	var r Release
	var err error
	var downloadLink string

	if args.Lastest {
		r, err = fetchLatestRelease(args.User, args.Project)
	} else {
		r, err = fetchLatestStableRelease(args.User, args.Project)
	}
	if err != nil {
		fmt.Printf("Failed to fetch releases for %s/%s: %s", args.User, args.Project, err)
		os.Exit(1)
	}

	if args.Kernel == "" {
		args.Kernel = GOOS
	}
	fmt.Printf("Your Kernel is: %s\n", args.Kernel)

	matchKernel := []string{}
	for _, asset := range r.Assets {
		if strings.Contains(strings.ToLower(asset.Name), strings.ToLower(args.Kernel)) {
			matchKernel = append(matchKernel, asset.URL)
		}
	}

	if args.Architecture == "auto" {
		args.Architecture = GOARCH
	}
	// Match Architecture
	if args.Architecture != "" {
		foundMatch := false
		for _, link := range matchKernel {
			if strings.Contains(strings.ToLower(link), strings.ToLower(args.Architecture)) {
				downloadLink = link
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			fmt.Println("Can not found release with your Architecture.")
			os.Exit(0)
		}
	}
	fmt.Printf("Your Architecture is: %s\n", GOARCH)

	fmt.Println(downloadLink)

	if args.Download {
		fmt.Println("Downloading")
		err = downloadFile(downloadLink)
		if err != nil {
			panic(err)
		}
		fmt.Println("Download Finish")
	}
}

func downloadFile(downloadUrl string) error {
	u, err := url.Parse(downloadUrl)
	if err != nil {
		return err
	}
	path := u.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	// Get the data
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
