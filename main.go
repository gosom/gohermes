package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
)

//go:embed templates/*.tmpl
var rootFs embed.FS

type appConfig struct {
	FolderName  string
	PackageName string
}

func main() {
	var (
		err error
		cfg appConfig
	)
	cfg.FolderName, err = stringPrompt("Enter folder name (no spaces)", "")
	if err != nil {
		log.Panic(err)
	}
	cfg.PackageName, err = stringPrompt("Enter the go package name (example: github.com/username/simple-api)", "")
	if err != nil {
		log.Panic(err)
	}

	if err := setUpDirectories(cfg); err != nil {
		log.Panic(err)
	}
	if err := setUpTemplates(cfg); err != nil {
		log.Panic(err)
	}

	fmt.Printf("\n🎉 Congratulations! Your new application is ready.")
	fmt.Printf("\nTo begin execute the following:\n\n")
	fmt.Printf("   cd %s\n", cfg.FolderName)
	fmt.Printf("   go run .\n")
}

func setUpDirectories(cfg appConfig) error {
	if err := os.Mkdir(cfg.FolderName, 0755); err != nil {
		return err
	}

	if err := os.Chdir(cfg.FolderName); err != nil {
		return err
	}

	subdirs := []string{
		"migrations",
		"models",
		"modelsext",
		"routes",
		"services",
		"user",
	}

	for _, dirname := range subdirs {
		if err := os.MkdirAll(dirname, 0755); err != nil {
			return err
		}
	}
	return nil
}

func setUpTemplates(cfg appConfig) error {
	tpl, err := template.ParseFS(
		rootFs, "templates/*.tmpl",
	)
	if err != nil {
		return err
	}
	templates := tpl.Templates()
	for i := range templates {
		err := func(t *template.Template) error {
			name := t.Name()
			var folder string
			var fname string
			if strings.HasSuffix(name, ".sql.tmpl") {
				folder = "./migrations"
				fname = name[len(folder) : len(name)-5]
			} else if name == "env.tmpl" {
				folder = "./"
				fname = ".env"
			} else {
				parts := strings.Split(name, "_")
				if len(parts) == 1 {
					folder = "."
				} else {
					folder = "./" + strings.Join(parts[:len(parts)-1], "/")
				}
				fname = parts[len(parts)-1]
				fname = fname[:len(fname)-5]
			}
			outputPath := strings.Join([]string{folder, fname}, "/")
			fp, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer fp.Close()
			return tpl.ExecuteTemplate(fp, name, cfg)

		}(templates[i])
		if err != nil {
			return err
		}
	}
	err = goModules(cfg)
	return err
}

func stringPrompt(label, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	return prompt.Run()
}

func goModules(cfg appConfig) error {
	cmd := exec.Command("go", "mod", "init", cfg.PackageName)
	if err := runCommand(cmd); err != nil {
		return err
	}
	cmd = exec.Command("go", "mod", "tidy")
	if err := runCommand(cmd); err != nil {
		return err
	}
	return nil
}

func runCommand(cmd *exec.Cmd) error {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	slurp, _ := io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
