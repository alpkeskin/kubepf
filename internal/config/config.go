/*
Copyright © 2023 github.com/alpkeskin
*/
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Projects []Project
}

type Project struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Services  []Service
}

type Service struct {
	Name       string `yaml:"name"`
	Localport  int    `yaml:"local_port"`
	Targetport int    `yaml:"target_port"`
}

var (
	Cfg *Config
)

func New() *Config {
	return &Config{}
}

func (c *Config) Parse() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(homeDir, ".kubepf")

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	var cFile Config
	err = yaml.Unmarshal(data, &cFile)
	if err != nil {
		return err
	}

	Cfg = &cFile

	return nil
}

func (c *Config) Init() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(homeDir, ".kubepf")

	if c.Exists() {
		return fmt.Errorf(".kubepf config file already exists")
	}

	// Create dummy config file
	cFile := Config{
		Projects: []Project{
			{
				Name:      "default",
				Namespace: "default",
				Services: []Service{
					{
						Name:       "service",
						Localport:  8080,
						Targetport: 8080,
					},
				},
			},
		},
	}

	data, err := yaml.Marshal(cFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println(color.GreenString(".kubepf config file created successfully in your home directory!\nNow you can edit this file and use kubepf commands."))
	return nil
}

func (c *Config) Exists() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	configFilePath := filepath.Join(homeDir, ".kubepf")
	_, err = os.Stat(configFilePath)

	return !os.IsNotExist(err)
}

func (c *Config) PrintList() {
	for _, project := range Cfg.Projects {
		fmt.Println(color.GreenString("* "), project.Name)
		for _, service := range project.Services {
			fmt.Println(color.YellowString("  - "), service.Name, color.CyanString("%d:%d", service.Localport, service.Targetport))
		}
		fmt.Println()
	}
}
