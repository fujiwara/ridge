package ridge

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Init struct {
		Name        string `help:"Name of the project" required:""`
		Description string `help:"Description of the project" default:""`
	} `cmd:"" help:"Initialize a new Ridge project"`

	Build struct {
	} `cmd:"" help:"Build the current Ridge project"`

	Dev struct {
		Port int `help:"Port to run the server on" default:"8080"`
	} `cmd:"" help:"Run the current Ridge project in development mode on localhost"`
}

func RunCLI(ctx context.Context) error {
	cli := &CLI{}
	kongCtx := kong.Parse(cli)
	switch kongCtx.Command() {
	case "init":
		return cli.RunInit(ctx)
	case "build":
		return cli.RunBuild(ctx)
	case "dev":
		return cli.RunDev(ctx)
	default:
	}
	return fmt.Errorf("unknown command: %s", kongCtx.Command())
}

func (cli *CLI) RunInit(ctx context.Context) error {
	src := `
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fujiwara/ridge"
)

func main() {
	var mux = http.NewServeMux()
	mux.HandleFunc("/", handleHello)
	ridge.Run(os.Getenv("RIDGE_ADDR"), "/", mux)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello %s\n", r.FormValue("name"))
}
`

	if _, err := os.Stat("main.go"); err == nil {
		return fmt.Errorf("main.go already exists")
	}
	log.Println("[info] creating main.go")
	if err := os.WriteFile("main.go", []byte(src), 0644); err != nil {
		return err
	}
	argss := [][]string{
		{"fmt", "main.go"},
		{"mod", "init"},
		{"mod", "tidy"},
		{"get", "."},
	}
	for _, args := range argss {
		log.Println("[info] running: go", strings.Join(args, " "))
		c := exec.CommandContext(ctx, "go", args...)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to run go %s: %w", strings.Join(args, " "), err)
		}
	}
	return nil
}

func (cli *CLI) RunBuild(ctx context.Context) error {
	args := []string{"build", "-o", "bootstrap", "main.go"}
	log.Println("[info] running: go", strings.Join(args, " "))
	c := exec.CommandContext(ctx, "go", args...)
	envs := os.Environ()
	for _, env := range envs {
		if strings.HasPrefix(env, "GOOS=") || strings.HasPrefix(env, "GOARCH=") || strings.HasPrefix(env, "CGO_ENABLED=") {
			continue
		}
		envs = append(envs, env)
	}
	envs = append(envs, "GOOS=linux")
	envs = append(envs, "GOARCH=amd64")
	envs = append(envs, "CGO_ENABLED=0")
	c.Env = envs
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

func (cli *CLI) RunDev(ctx context.Context) error {
	args := []string{"run", "main.go"}
	log.Println("[info] running: go", strings.Join(args, " "))
	c := exec.CommandContext(ctx, "go", args...)
	c.Env = os.Environ()
	c.Env = append(c.Env, "RIDGE_ADDR="+fmt.Sprintf(":%d", cli.Dev.Port))
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}
