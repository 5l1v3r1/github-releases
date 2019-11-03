package main

import (
	"fmt"
	"os"
	"strings"

	arg "github.com/alexflint/go-arg"
)

var args struct {
	User    string `arg:"required"`
	Project string `arg:"required"`
	// Version string `arg:""`
	Distro  string `arg:"" default:"linux" help:"linux/darwin/window"`
	Lastest bool   `arg:"" help:"Get last version"`
}

func main() {
	arg.MustParse(&args)

	var r Release
	var err error
	if args.Lastest {
		r, err = fetchLatestRelease(args.User, args.Project)
	} else {
		r, err = fetchLatestStableRelease(args.User, args.Project)
	}
	if err != nil {
		fmt.Printf("Failed to fetch releases for %s/%s: %s", args.User, args.Project, err)
		os.Exit(1)
	}

	for i := 0; i < len(r.Assets); i++ {
		if strings.Contains(strings.ToLower(r.Assets[i].Name), strings.ToLower(args.Distro)) {
			fmt.Println(r.Assets[i].URL)
			break
		}
	}
	// fmt.Println(r.Assets)
	// args.Version = strings.TrimPrefix(args.Version, "v")
	// r.Name = strings.TrimPrefix(r.Name, "v")
	// if r.Name != args.Version {
	// 	fmt.Printf("Not latest. Your Version %s - Latest: %s\n", args.Version, r.Name)
	// 	if len(r.Assets) > 0 && args.URL {
	// 		fmt.Printf("URL: %s\n", r.Assets[0].URL)
	// 	}
	// 	os.Exit(1)
	// }
	// os.Exit(0)
}
