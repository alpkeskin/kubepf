/*
Copyright © 2023 github.com/alpkeskin
*/
package lighter

import (
	"fmt"
	"os"

	"github.com/alpkeskin/kubepf/internal/config"
	"github.com/alpkeskin/kubepf/pkg/shell"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kubepf [command]",
	Short:   "\nFaster Way to Use Kubectl Port Forwarding",
	Long:    "\nFaster Way to Use Kubectl Port Forwarding",
	Run:     magic,
	Version: "v1.0.0",
}

func Fire() {
	rootCmd.Example = "\n  kubepf [project-name]\n  kubepf list\n  kubepf active\n  kubepf kill [project-name]\n  kubepf kill [service-name]"
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func magic(cmd *cobra.Command, args []string) {
	// Create config instance
	conf := config.New()

	// Check config file exists
	if !conf.Exists() {
		fmt.Println(color.RedString("Create .kubepf config file in your home directory!"))
		os.Exit(1)
	}

	// Parse config file
	err := conf.Parse()
	if err != nil {
		panic(err)
	}

	// Check args. if empty, print help
	if len(args) == 0 {
		cmd.Help()
		return
	}

	input := args[0]

	// Check input. If list, print config file
	if input == "list" {
		conf.PrintList()
		return
	}

	shell := shell.New()

	// Check input. If active, print active port-forwards
	if input == "active" {
		shell.Active(*config.Cfg)
		return
	}

	// Check input. If kill, kill port-forward
	if input == "kill" {
		if len(args) < 2 {
			fmt.Println(color.RedString("Please enter a service name or project name"))
			os.Exit(1)
		}

		for _, project := range config.Cfg.Projects {
			if project.Name == args[1] {
				for _, service := range project.Services {
					err := shell.Kill(service.Name)
					if err != nil {
						panic(err)
					}
				}
				return
			}
		}

		err := shell.Kill(args[1])
		if err != nil {
			panic(err)
		}

		return
	}

	// set commands
	commands := setCommands(input)

	// execute commands
	for _, command := range commands {
		err := shell.Exec(command)
		if err != nil {
			panic(err)
		}
		fmt.Println(color.GreenString("✓"), command, color.CyanString("[started]"))
	}
}

func setCommands(input string) []string {
	commands := []string{}
	for _, project := range config.Cfg.Projects {
		if project.Name == input {
			for _, service := range project.Services {
				command := fmt.Sprintf("kubectl port-forward --namespace %s svc/%s %d:%d", project.Namespace, service.Name, service.Localport, service.Targetport)
				commands = append(commands, command)
			}
		}
	}
	return commands
}
