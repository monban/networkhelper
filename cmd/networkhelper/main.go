package main

import (
	"os"
	"os/exec"

	_ "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/erikgeiser/promptkit/confirmation"
)

func main() {
	c := confirmation.New("Connect DHCP?", confirmation.Yes)
	dhcp, err := c.RunPrompt()
	if err != nil {
		log.Fatal(err)
	}
	if dhcp {
		x("udhcpc")
	}
}

func x(cmd string) {
	if os.Getenv("USER") == "root" {
		exec.Command(cmd).Run()
	} else {
		log.Infof("SIMULATION: %s", cmd)
	}
}
