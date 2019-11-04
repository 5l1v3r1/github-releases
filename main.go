package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/matishsiao/goInfo"
)

var args struct {
	User     string `arg:"required"`
	Project  string `arg:"required"`
	Kernel   string `arg:"" default:"" help:"linux/darwin/window (Default is your running Kernel)"`
	Platform string `arg:"" default:"" help:"x86_64/i386"`
	Lastest  bool   `arg:"" help:"Get last release (Default is Latest Stable Release)"`
	Download bool   `arg:"" help:"Download file"`
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

	gi := goInfo.GetInfo()
	if args.Kernel == "" {
		args.Kernel = gi.Kernel
	}

	for i := 0; i < len(r.Assets); i++ {
		if strings.Contains(strings.ToLower(r.Assets[i].Name), strings.ToLower(args.Kernel)) {
			downloadLink = r.Assets[i].URL
			break
		}
	}

	if args.Platform != "" {
		// Match Platform
		if strings.Contains(strings.ToLower(downloadLink), strings.ToLower(args.Platform)) {
			fmt.Println(downloadLink)
			os.Exit(0)
		}
	}
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
