package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	wit "github.com/patrickhuber/go-wasm/wit/parse"
)

const witFileExtension = ".wit"

func main() {
	flagWIT := flag.String("wit", "wit", "WIT directory to process")
	flagOut := flag.String("out", "gen", "output directory to write generated files into")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	goPackageBasename := args[0]
	outputRoot := *flagOut
	rootFS := os.DirFS(*flagWIT)

	if err := fs.WalkDir(rootFS, ".", func(path string, dir os.DirEntry, err error) error {
		// Skip files that do not have the right extension.
		if !strings.HasSuffix(path, witFileExtension) {
			return nil
		}

		fmt.Printf("---> Loading %s\n", path)

		// Read the file.
		content, err := fs.ReadFile(rootFS, path)
		if err != nil {
			return err
		}

		// Parse the content.
		witAST, err := wit.Parse(string(content))
		if err != nil {
			return err
		}

		// Generate output packages.
		return generate(filepath.SplitList(goPackageBasename), path, outputRoot, witAST)
	}); err != nil {
		log.Fatalln(err)
	}
}
