package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var (
	sysLog *log.Logger
	errLog *log.Logger
)

// Delete path from the filesystem
func kaboom(path string) error {
	fd, err := os.Stat(path)
	if err != nil {
		errLog.Fatalf("error getting file stats : %v", err)
	}

	switch mode := fd.Mode(); {
	case mode.IsDir():
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	case mode.IsRegular():
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	sysLog = log.New(os.Stdout, "fs-bazooka: ", log.Ldate|log.Ltime|log.Lshortfile)
	errLog = log.New(os.Stdout, "fs-bazooka: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	var (
		paths  []string
		dir    = flag.String("path", "", "Nuking will happen in this directory path")
		upshot = flag.Float64("upshot", 6.0, "Outcome of events")
	)
	flag.Parse()

	prob := float64((1 / *upshot) * 100)
	if math.IsInf(prob, 1) {
		errLog.Fatalf("upshot cannot be 0")
	}

	sysLog.Printf("Probability : %f", prob)
	bombSite := *dir
	err := filepath.Walk(bombSite, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errLog.Fatalf("error accessing a path %q: %v\n", bombSite, err)
			return err
		}
		paths = append(paths, path)

		return nil
	})

	if err != nil {
		errLog.Fatalf("error walking the path %q: %v\n", bombSite, err)
	}

	// Determine the unlucky path
	source1 := rand.NewSource(time.Now().UnixNano())
	rand1 := rand.New(source1)
	unlucky := rand1.Intn(len(paths))

	// Determine to take the shot or not
	source2 := rand.NewSource(time.Now().UnixNano())
	rand2 := rand.New(source2)
	shot := rand2.Intn(int(*upshot))

	if shot == 0 {
		if err := kaboom(paths[unlucky]); err != nil {
			errLog.Fatalf("error removing the file / directory : %v", err)
		}
		sysLog.Printf("unlucky path %s", paths[unlucky])
	} else {
		sysLog.Printf("lucky path %s", paths[unlucky])
	}
}
