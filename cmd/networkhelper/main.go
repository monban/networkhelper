package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	c := confirmation.New("Connect DHCP?", confirmation.Yes)
	dhcp, err := c.RunPrompt()
	if err != nil {
		log.Fatal(err)
	}
	if dhcp {
		x("udhcpc")
	} else {
		var ip, mask, gw string
		ip, _ = ipinput("IP address:", "").RunPrompt()
		ipBytes := net.ParseIP(ip)
		log.Info("len", "ipBytes", len(ipBytes))
		if ipBytes.To4() == nil { //ipv6 town
			gw, _ = ipinput("Default gateway:", "").RunPrompt()
			x("ifconfig", "eth0", ip)
		} else { // ipv4
			maskBytes := ipBytes.DefaultMask()
			ipBytes[len(ipBytes)-1]++
			mask, _ = maskinput(maskBytes).RunPrompt()
			gwInitial := ip[:strings.LastIndex(ip, ".")+1]
			gw, _ = ipinput("Default gateway:", gwInitial).RunPrompt()
			x("ifconfig", "eth0", ip, "netmask", mask)
		}
		x("route", "add", "default", "gw", gw)
	}
}

func x(cmd string, args ...string) {
	if os.Getenv("USER") == "root" {
		exec.Command(cmd, args...).Run()
	} else {
		log.Infof("SIMULATION: %s %s", cmd, args)
	}
}

func ipinput(prompt string, dlft string) *textinput.TextInput {
	input := textinput.New(prompt)
	input.InitialValue = dlft
	input.Placeholder = "192.168.1.10"
	input.AutoComplete = textinput.AutoCompleteFromSliceWithDefault([]string{
		"10.0.0.1",
		"10.0.0.2",
		"127.0.0.1",
		"fe80::1",
	}, input.Placeholder)

	input.Validate = func(input string) error {
		if net.ParseIP(input) == nil {
			return fmt.Errorf("invalid IP address")
		}
		return nil
	}
	return input
}

func maskinput(mask net.IPMask) *textinput.TextInput {
	input := textinput.New("Network mask:")
	input.InitialValue = net.IP(mask).String()
	return input
}
