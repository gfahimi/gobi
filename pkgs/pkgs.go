package pkgs

import (
	"os"
	"path/filepath"
	"strings"
)

// PkgSet defines a package set.
type PkgSet map[string]struct{}

// List returns a list of packages.
func (p PkgSet) List() []string {
	s := make([]string, 0, len(p))
	for n := range p {
		s = append(s, n)
	}
	return s
}

// Add a package to the package set.
func (p PkgSet) Add(pkg string) {
	p[pkg] = struct{}{}
}

// Packages returns all go packages in the project dir. It ignores the skipped packages.
func Packages(projectDir string, skipPkgs PkgSet) (PkgSet, error) {
	skipPkgs["vendor"] = struct{}{}
	pkgs := PkgSet{}
	err := filepath.Walk(projectDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				name := info.Name()
				if _, ok := skipPkgs[name]; ok {
					return filepath.SkipDir
				}
				if name != "." && strings.HasPrefix(name, ".") {
					return filepath.SkipDir
				}
			}
			if strings.HasSuffix(info.Name(), ".go") {
				// Add this directory to the packages
				p := strings.SplitN(path, "/", 2)
				if len(p) == 1 {
					pkgs.Add("")
				} else {
					pkgs.Add(p[0])
				}
			}
			return nil
		})

	if err != nil {
		return nil, err
	}
	return pkgs, nil
}
