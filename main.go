package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/manifoldco/promptui"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     pkg/scheduler/schedulerrpc/scheduler.proto

//go:embed templates/*.tmpl
var rootFs embed.FS

type appConfig struct {
	AppName     string
	PackageName string

	ServerAddress string

	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string

	TokenSecret string

	UseScheduler string
	Debug        string

	DockerNetwork string
	DockerDbHost  string

	SchedulerServerAddress string

	SchedulerDbHost     string
	SchedulerDbPort     string
	SchedulerDbName     string
	SchedulerDbUser     string
	SchedulerDbPassword string
	SchedulerSecret     string
	SchedulerApiKey     string
	SchedulerDebug      string
}

func main() {
	cfg, err := readConfig()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Setting up project skeleton")

	if err := setUpDirectories(cfg); err != nil {
		log.Panic(err)
	}
	if err := setUpTemplates(cfg); err != nil {
		log.Panic(err)
	}

	log.Println("setting up database (it will try to ping db for 3 minutes or until docker-compose is ready)")
	if err := setUpDb(cfg); err != nil {
		log.Panic(err)
	}

	log.Println("Initializing go modules")
	if err := goModules(cfg); err != nil {
		log.Panic(err)
	}

	log.Println("generating models from db")
	if err := generateModels(cfg); err != nil {
		log.Panic(err)
	}

	fmt.Printf("\n🎉 Congratulations! Your new application is ready.")
	fmt.Printf("\nTo begin see the makefile:\n\n")
	fmt.Printf("   cd %s\n", cfg.AppName)
	fmt.Printf("   make\n")
}

func readConfig() (appConfig, error) {
	type entry struct {
		message        string
		defaultValue   string
		storeTo        *string
		acceptedValues []string
	}
	var cfg appConfig
	params := []entry{
		{message: "Enter app name (no spaces)", storeTo: &cfg.AppName},
		{message: "Enter go package name", storeTo: &cfg.PackageName},
		{message: "Enter listen address", defaultValue: ":8080", storeTo: &cfg.ServerAddress},

		{message: "Enter database host", defaultValue: "localhost", storeTo: &cfg.DbHost},
		{message: "Enter database port", defaultValue: "5432", storeTo: &cfg.DbPort},
		{message: "Enter database name", storeTo: &cfg.DbName},
		{message: "Enter database user", storeTo: &cfg.DbUser},
		{message: "Enter database password", storeTo: &cfg.DbPassword},
		{message: "Enter secret", storeTo: &cfg.TokenSecret},
		{message: "Enable debug? (true || false)", defaultValue: "true", storeTo: &cfg.Debug,
			acceptedValues: []string{"true", "false"}},
		{message: "Do you need a scheduler? (true || false)", defaultValue: "true", storeTo: &cfg.UseScheduler},
	}

	for _, p := range params {
		v, err := stringPrompt(p.message, p.defaultValue)
		if err != nil {
			return cfg, err
		}
		if len(p.acceptedValues) > 0 {
			isAccepted := false
			for _, accepted := range p.acceptedValues {
				if accepted == v {
					isAccepted = true
					break
				}
			}
			if !isAccepted {
				return cfg, fmt.Errorf("value %s not accepted", v)
			}
		}
		*p.storeTo = v
	}

	if cfg.UseScheduler == "true" {
		schedulerParams := []entry{
			{message: "Enter scheduler listen address", defaultValue: ":50051", storeTo: &cfg.SchedulerServerAddress},
			{message: "Enter scheduler database host", defaultValue: cfg.DbHost, storeTo: &cfg.SchedulerDbHost},
			{message: "Enter scheduler database port", defaultValue: cfg.DbPort, storeTo: &cfg.SchedulerDbPort},
			{message: "Enter scheduler database name", defaultValue: cfg.DbName, storeTo: &cfg.SchedulerDbName},
			{message: "Enter scheduler database user", defaultValue: cfg.DbUser, storeTo: &cfg.SchedulerDbUser},
			{message: "Enter scheduler database password", defaultValue: cfg.DbPassword, storeTo: &cfg.SchedulerDbPassword},
			{message: "Enter apiKey", storeTo: &cfg.SchedulerApiKey},
			{message: "Enter secret", storeTo: &cfg.SchedulerSecret},
			{message: "Enable debug? (true || false)", defaultValue: "true", storeTo: &cfg.SchedulerDebug,
				acceptedValues: []string{"true", "false"}},
		}
		for _, p := range schedulerParams {
			v, err := stringPrompt(p.message, p.defaultValue)
			if err != nil {
				return cfg, err
			}
			if len(p.acceptedValues) > 0 {
				isAccepted := false
				for _, accepted := range p.acceptedValues {
					if accepted == v {
						isAccepted = true
						break
					}
				}
				if !isAccepted {
					return cfg, fmt.Errorf("value %s not accepted", v)
				}
			}
			*p.storeTo = v
		}
	}

	cfg.DockerNetwork = "network_" + strings.ReplaceAll(cfg.AppName, "/", "_")
	cfg.DockerDbHost = "db"
	return cfg, nil
}

func setUpDirectories(cfg appConfig) error {
	if err := os.Mkdir(cfg.AppName, 0755); err != nil {
		return err
	}

	if err := os.Chdir(cfg.AppName); err != nil {
		return err
	}

	subdirs := []string{
		"commands",
		"migrations",
		"models",
		"modelsext",
		"routes",
		"services",
		"user",
	}
	if cfg.UseScheduler == "true" {
		subdirs = append(subdirs, "scheduler")
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
	return nil
}

func setUpDb(cfg appConfig) error {
	if err := startDockerDb(cfg); err != nil {
		return err
	}
	cmd := exec.Command("make", "migrate-up")
	return runCommand(cmd)
}

func generateModels(cfg appConfig) error {
	cmd := exec.Command("make", "generate-models")
	return runCommand(cmd)
}

func startDockerDb(cfg appConfig) error {
	cmd := exec.Command("make", "db-start")
	if err := runCommand(cmd); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(5 * time.Second)
			if err := testDbConnection(ctx, cfg); err == nil {
				return nil
			} else {
				log.Println("Error pinging db: " + err.Error())
			}
		}
	}
}

func testDbConnection(ctx context.Context, cfg appConfig) error {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbUser, cfg.DbPassword)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.PingContext(ctx)
}

func stringPrompt(label, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	return prompt.Run()
}

func goModules(cfg appConfig) error {
	cmd := exec.Command("go", "mod", "download", cfg.PackageName)
	if err := runCommand(cmd); err != nil {
		return err
	}
	return nil
}

func runCommand(cmd *exec.Cmd) error {
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
