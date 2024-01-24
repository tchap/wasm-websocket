package main

import (
	"flag"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"strings"

	wit "github.com/patrickhuber/go-wasm/wit/parse"
	"github.com/tchap/wamcloud-websocket/bin/provider-with-bindgen/internal/generator"
)

const witFileExtension = ".wit"

func main() {
	flagIn := flag.String("in", "wit", "input WIT directory to process")
	flag.Parse()

	rootFS := os.DirFS(*flagIn)

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
		pkgAST, err := wit.Parse(string(content))
		if err != nil {
			return err
		}

		// Generate output packages.
		out := generator.FormatFile(generator.BuildFile(pkgAST))
		outFmt, err := format.Source(out)
		if err != nil {
			fmt.Println(string(out))
			return err
		}

		fmt.Println(string(outFmt))
		return nil
	}); err != nil {
		log.Fatalln(err)
	}
}
