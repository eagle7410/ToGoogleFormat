package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

type match struct {
	re *regexp.Regexp
}

func (i *match) isGoFile(file os.FileInfo) bool {
	return i.re.MatchString(file.Name())
}

func (i *match) init() {
	i.re = regexp.MustCompile(`(.*)\.go`)
}

var m match

func init() {
	m.init()
}

func main() {
	base := "."
	dir := ""
	workDir(&base, &dir)
}

func workDir(base, dir *string) {
	wd := *base

	if len(*dir) > 0 {
		wd += "/" + *dir
	}

	files, err := ioutil.ReadDir(wd)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {
			name := file.Name()
			workDir(&wd, &name)
			continue
		}

		if m.isGoFile(file) {
			filePath := wd + "/" + file.Name()

			if err := cmdFmt(&filePath); err != nil {
				fmt.Printf("[0;31mError cmd for file: %v %v[39m\n", filePath, err)
			} else {
				fmt.Printf("file: %v is Google format Ok \n", filePath)
			}
		}
	}
}

func cmdFmt(filePath *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "fmt", *filePath)
	_, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("Deadline exceeded for go fmt " + *filePath)
	}

	return err
}
