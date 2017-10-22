package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	pwd = "./"
)

func init() {
	build.Default.UseAllFiles = true
}

func usage() {
	log.Printf("Usage of %s:\n", os.Args[0])
	log.Printf("\nunimport # runs on package in current directory\n")
	log.Printf("\nunimport [packages]\n")
	//TODO add back when flags are supporrted
	// log.Printf("Flags:\n")
	// flag.PrintDefaults()
}

type returnsVisitor struct {
	f *token.FileSet
}

func main() {
	// Remove log timestamp
	log.SetFlags(0)

	flag.Usage = usage
	flag.Parse()

	if err := checkImports(flag.Args()); err != nil {
		log.Println(err)
	}

}

func checkImports(args []string) error {

	fset := token.NewFileSet()
	files, err := parseInput(args, fset)
	if err != nil {
		return fmt.Errorf("could not parse input %v", err)
	}

	retVis := &returnsVisitor{
		f: fset,
	}

	for _, f := range files {
		ast.Walk(retVis, f)
	}

	return nil
}

func parseInput(args []string, fset *token.FileSet) ([]*ast.File, error) {
	var directoryList []string
	var fileMode bool
	files := make([]*ast.File, 0)

	if len(args) == 0 {
		directoryList = append(directoryList, pwd)
	} else {
		for _, arg := range args {
			if strings.HasSuffix(arg, "/...") && isDir(arg[:len(arg)-len("/...")]) {

				for _, dirname := range allPackagesInFS(arg) {
					directoryList = append(directoryList, dirname)
				}

			} else if isDir(arg) {
				directoryList = append(directoryList, arg)

			} else if exists(arg) {
				if strings.HasSuffix(arg, ".go") {
					fileMode = true
					f, err := parser.ParseFile(fset, arg, nil, 0)
					if err != nil {
						return nil, err
					}
					files = append(files, f)
				} else {
					return nil, fmt.Errorf("invalid file %v specified", arg)
				}
			} else {

				//TODO clean this up a bit
				imPaths := importPaths([]string{arg})
				for _, importPath := range imPaths {
					pkg, err := build.Import(importPath, ".", 0)
					if err != nil {
						return nil, err
					}
					var stringFiles []string
					stringFiles = append(stringFiles, pkg.GoFiles...)
					// files = append(files, pkg.CgoFiles...)
					stringFiles = append(stringFiles, pkg.TestGoFiles...)
					if pkg.Dir != "." {
						for i, f := range stringFiles {
							stringFiles[i] = filepath.Join(pkg.Dir, f)
						}
					}

					fileMode = true
					for _, stringFile := range stringFiles {
						f, err := parser.ParseFile(fset, stringFile, nil, 0)
						if err != nil {
							return nil, err
						}
						files = append(files, f)
					}

				}
			}
		}
	}

	// if we're not in file mode, then we need to grab each and every package in each directory
	// we can to grab all the files
	if !fileMode {
		for _, fpath := range directoryList {
			pkgs, err := parser.ParseDir(fset, fpath, nil, 0)
			if err != nil {
				return nil, err
			}

			for _, pkg := range pkgs {
				for _, f := range pkg.Files {
					files = append(files, f)
				}
			}
		}
	}

	return files, nil
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (v *returnsVisitor) Visit(node ast.Node) ast.Visitor {

	file, ok := node.(*ast.File)
	if !ok {
		return v
	}

	type importAlias struct {
		importSpec *ast.ImportSpec
		index      int
	}

	var importAliases []importAlias
	for i, pkgImport := range file.Imports {
		if pkgImport.Name != nil && pkgImport.Name.Name != "_" {
			alias := importAlias{
				importSpec: pkgImport,
				index:      i,
			}
			importAliases = append(importAliases, alias)
		}
	}

	switch len(importAliases) {
	case 0:

	default:
		// verify that each alias is needed by making a second pass through the imports
		for _, importAlias := range importAliases {
			var aliasNeeded bool
			for i, pkgImport := range file.Imports {
				// Since we know the index of the import alias in file.Imports from our first pass, we can save a string comparison
				if i == importAlias.index {
					continue
				}
				if pkgImport.Path != nil && strings.Replace(path.Base(pkgImport.Path.Value), `"`, "", -1) == strings.Replace(path.Base(importAlias.importSpec.Path.Value), `"`, "", -1) {
					// this alias is needed, continue
					aliasNeeded = true
					break
				}

			}
			if !aliasNeeded {
				file := v.f.File(importAlias.importSpec.Pos())
				lineNumber := file.Position(importAlias.importSpec.Pos()).Line
				// dot imports inside of tests are okay
				if importAliases[0].importSpec.Name.Name == "." && strings.HasSuffix(file.Name(), "_test.go") {
					continue
				}
				// If the alias path contains a dash or dot, it's likely importing a specific revision - ignore these
				if strings.Contains(importAliases[0].importSpec.Path.Value, "-") || strings.Contains(importAliases[0].importSpec.Path.Value, ".") {
					continue
				}
				log.Printf("%v:%v unnecessary import alias %v\n", file.Name(), lineNumber, importAlias.importSpec.Name.Name)
			}
		}
	}

	return v
}
