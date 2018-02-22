package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/anuvu/gobi/lint"
	"github.com/anuvu/gobi/pkgs"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// Builder provides build capability with common build targets like lint,
// test and binary builds.
type Builder struct {
	ctx         *cli.Context
	pkgs        pkgs.PkgSet
	ExcludePkgs []string     `yaml:"exclude"`
	Linter      *lint.Linter `yaml:"lint"`
}

// New returns a new instance of builder
func New(ctx *cli.Context) (*Builder, error) {
	conf := ctx.String("config-file")
	f, e := os.Open(conf)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	d, e := ioutil.ReadAll(f)
	if e != nil {
		return nil, e
	}

	b := &Builder{}
	if err := yaml.Unmarshal(d, b); err != nil {
		return nil, err
	}

	edirs := pkgs.PkgSet{}
	for _, edir := range b.ExcludePkgs {
		edirs[edir] = struct{}{}
	}

	p, e := pkgs.Packages(ctx.String("project-dir"), edirs)
	if e != nil {
		return nil, e
	}
	b.pkgs = p
	return b, nil
}

// Execute runs all targets. Returns nil on success, else error.
func (b *Builder) Execute() error {
	startTime := time.Now()
	defer func() {
		fmt.Printf("Time Taken: %v\n", time.Now().Sub(startTime))
	}()
	fmt.Println("=== Build started ===")
	err := b.Linter.Execute(b.ctx, b.pkgs)
	if err != nil {
		return err
	}
	return nil
}

// Build executes all declared build targets in the dependency order.
func Build(c *cli.Context) error {
	builder, err := New(c)
	if err != nil {
		return err
	}
	return builder.Execute()
}
