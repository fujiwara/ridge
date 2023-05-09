package ridgecli

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/fatih/color"
	"github.com/fujiwara/logutils"
)

type CLI struct {
	LogLevel string `help:"Log level" default:"info" enum:"trace,debug,info,warn,error"`

	Init struct {
		Name        string `help:"Name of the project" required:""`
		Description string `help:"Description of the project" default:""`
		AccountID   string `help:"AWS Account ID"`
	} `cmd:"" help:"Initialize a new Ridge project"`

	Build struct {
	} `cmd:"" help:"Build the current Ridge project"`

	Dev struct {
		Port int `help:"Port to run the server on" default:"8080"`
	} `cmd:"" help:"Run the current Ridge project in development mode on localhost"`

	Deploy struct {
		DryRun bool `help:"Dry run" default:"false"`
		Build  bool `help:"Build before deploy" default:"true" negatable:""`
	} `cmd:"" help:"Deploy the current Ridge project to AWS Lambda"`

	awscfg aws.Config
}

func Run(ctx context.Context) error {
	cli := &CLI{}
	kongCtx := kong.Parse(cli)
	if err := cli.prepare(ctx); err != nil {
		return err
	}
	switch kongCtx.Command() {
	case "init":
		return cli.RunInit(ctx)
	case "build":
		return cli.RunBuild(ctx)
	case "dev":
		return cli.RunDev(ctx)
	case "deploy":
		return cli.RunDeploy(ctx)
	default:
	}
	return fmt.Errorf("unknown command: %s", kongCtx.Command())
}

type generateFileInfo struct {
	path       string
	src        []byte
	hook       func(context.Context, string) error
	isTemplate bool
}

//go:embed embed/main.go
var mainGoSrc []byte
var goFileHook = func(ctx context.Context, path string) error {
	argss := [][]string{
		{"fmt", path},
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

//go:embed embed/function.json
var functionJSONSrc []byte

//go:embed embed/.lambdaignore
var lambdaIgnoreSrc []byte

var generateFiles = []generateFileInfo{
	{
		path: "main.go",
		src:  mainGoSrc,
		hook: goFileHook,
	},
	{
		path:       "function.json",
		src:        functionJSONSrc,
		isTemplate: true,
	},
	{
		path: ".lambdaignore",
		src:  lambdaIgnoreSrc,
	},
}

func (cli *CLI) prepare(ctx context.Context) error {
	var err error
	cli.awscfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return err
	}

	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"trace", "debug", "info", "warn", "error"},
		ModifierFuncs: []logutils.ModifierFunc{
			logutils.Color(color.FgHiBlack), // trace
			logutils.Color(color.FgHiBlack), // debug
			nil,                             // info
			logutils.Color(color.FgYellow),  // warn
			logutils.Color(color.FgRed),     // error
		},
		MinLevel: logutils.LogLevel(cli.LogLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	return nil
}

func (cli *CLI) RunInit(ctx context.Context) error {
	log.Println("[info] initializing Ridge project")
	client := sts.NewFromConfig(cli.awscfg)
	out, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}
	cli.Init.AccountID = *out.Account
	log.Println("[info] AWS AccountID:", cli.Init.AccountID)
	for _, f := range generateFiles {
		if err := cli.generateFile(ctx, f); err != nil {
			return fmt.Errorf("failed to create %s: %w", f.path, err)
		}
	}
	return nil
}

func (cli *CLI) generateFile(ctx context.Context, info generateFileInfo) error {
	if _, err := os.Stat(info.path); err == nil {
		return fmt.Errorf("%s already exists", info.path)
	}
	log.Printf("[info] creating %s", info.path)

	f, err := os.Create(info.path)
	if err != nil {
		return err
	}
	if info.isTemplate {
		// template file
		tmpl := template.Must(template.New(info.path).Parse(string(info.src)))
		if err := tmpl.Execute(f, cli.Init); err != nil {
			return err
		}
	} else {
		// normal file
		if _, err := f.Write(info.src); err != nil {
			return err
		}
	}
	if info.hook == nil {
		return nil
	}
	return info.hook(ctx, info.path)
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

func (cli *CLI) RunDeploy(ctx context.Context) error {
	if cli.Deploy.Build {
		if err := cli.RunBuild(ctx); err != nil {
			return err
		}
	} else {
		log.Println("[info] skipping build")
	}
	args := []string{
		"deploy",
		"--log-level", cli.LogLevel,
	}
	if cli.Deploy.DryRun {
		args = append(args, "--dry-run")
	}
	c := exec.CommandContext(ctx, "lambroll", args...)
	c.Env = os.Environ()
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}
