package lint

import (
	"context"
	"os/exec"
	"time"

	"github.com/anuvu/gobi/pkgs"
	"github.com/anuvu/gobi/tasks"
	"github.com/urfave/cli"
)

// Linter provides the lint configuration for the build.
type Linter struct {
	Disabled    bool     `yaml:"disabled"`
	ExcludePkgs []string `yaml:"exclude"`
	ExtraLints  []string `yaml:"extraLints"`
	Timeout     string   `yaml:"timeout"`
	timeout     time.Duration
}

// Execute executes the lint target and returns true if it succeeded
// else false.
func (l *Linter) Execute(c *cli.Context, packages pkgs.PkgSet) error {
	if l != nil && l.Disabled {
		return nil
	}

	ts := tasks.New("lint", "./build", nil)
	exPkgs := pkgs.PkgSet{}
	if l != nil {
		for _, e := range l.ExcludePkgs {
			exPkgs.Add(e)
		}
	}

	for p := range packages {
		if _, ok := exPkgs[p]; !ok {
			ts.Add(p, l.makeCmd(p))
		}
	}

	e := ts.Run()
	ts.PrintStatus()
	return e
}

func (l *Linter) makeCmd(pkg string) *exec.Cmd {

	args := []string{
		"--disable-all",
		"--enable=golint",
		"--enable=vet",
		"--enable=gofmt",
	}

	if l != nil {
		for _, el := range l.ExtraLints {
			args = append(args, el)
		}
	}

	if pkg != "" {
		pkg = "./" + pkg + "/..."
	} else {
		pkg = "."
	}
	args = append(args, pkg)

	return exec.CommandContext(context.Background(), "gometalinter", args...)
}
