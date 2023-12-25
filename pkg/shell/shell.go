/*
Copyright © 2023 github.com/alpkeskin
*/
package shell

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/alpkeskin/kubepf/internal/config"
	"github.com/fatih/color"
)

type Shell struct {
	ActiveList []Active
}

type Active struct {
	Namespace string
	Service   string
}

func New() *Shell {
	return &Shell{}
}

func (s *Shell) Exec(cmd string) error {
	ctx := context.Background()
	timeout := time.Duration(900) * time.Second
	execCtx, _ := context.WithTimeout(ctx, timeout)

	cmdObj := exec.CommandContext(execCtx, "bash", "-c", cmd)

	if err := cmdObj.Start(); err != nil {
		return err
	}

	return nil
}

func (s *Shell) Active(cfg config.Config) error {
	cmd := "ps aux -o pid -o command | grep kubectl | grep port-forward"
	ctx := context.Background()
	timeout := time.Duration(900) * time.Second
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	out, err := exec.
		CommandContext(execCtx, "bash", "-c", cmd).
		Output()

	if err != nil {
		return err
	}

	regex := `--namespace\s+([^[:space:]]+).*svc\/([^[:space:]]+)`
	re := regexp.MustCompile(regex)

	for _, line := range re.FindAllStringSubmatch(string(out), -1) {
		s.ActiveList = append(s.ActiveList, Active{
			Namespace: line[1],
			Service:   line[2],
		})
	}

	for _, project := range cfg.Projects {
		fmt.Println(color.GreenString("* "), project.Name)
		for _, service := range project.Services {
			for _, active := range s.ActiveList {
				if project.Namespace == active.Namespace && service.Name == active.Service {
					fmt.Println(color.YellowString("  - "), service.Name, color.CyanString("%d:%d", service.Localport, service.Targetport), color.GreenString("\u2714"))
				}
			}
		}
		fmt.Println()
	}

	return nil
}

func (s *Shell) Kill(service string) error {
	cmd := fmt.Sprintf("ps aux -o pid -o command | grep kubectl | grep port-forward | grep %s | awk '{print $2}'", service)
	ctx := context.Background()
	timeout := time.Duration(900) * time.Second
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	out, err := exec.
		CommandContext(execCtx, "bash", "-c", cmd).
		Output()

	if err != nil {
		return err
	}

	regex := `([0-9]+)`
	re := regexp.MustCompile(regex)

	for _, line := range re.FindAllStringSubmatch(string(out), -1) {
		cmd := fmt.Sprintf("kill -9 %s", line[1])
		err := s.Exec(cmd)
		if err != nil {
			return err
		}
	}

	fmt.Println(color.GreenString("✓"), service, color.RedString("[killed]"))
	return nil
}
