package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/textinput"
	netmask "github.com/monban/bubble-netmask"
)

func main() {
	f, err := os.Create("menu.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	// n := textinput.NewModel(textinput.New("hello"))
	var m Menu = New()
	var p *tea.Program = tea.NewProgram(m, tea.WithAltScreen())
	p.Run()
	os.Exit(0)

	c := confirmation.New("Connect DHCP?", confirmation.Yes)
	dhcp, err := c.RunPrompt()
	if err != nil {
		log.Fatal(err)
	}
	if dhcp {
		x("udhcpc")
		x("ifconfig")
		x("route")
	} else {
		var ip, gw string
		var dns []string
		ip, _ = ipinput("IP address:", "").RunPrompt()
		ipBytes := net.ParseIP(ip)
		if ipBytes.To4() == nil { //ipv6 town
			gw, _ = ipinput("Default gateway:", "").RunPrompt()
			x("ifconfig", "eth0", ip)
		} else { // ipv4
			ipBytes[len(ipBytes)-1]++
			mask, err := netmask.New(ipBytes.String()).Run()
			if err != nil {
				log.Fatal(err)
			}
			gwInitial := ip[:strings.LastIndex(ip, ".")+1]
			gw, _ = ipinput("Default gateway:", gwInitial).RunPrompt()
			dns = getDNSServers()
			x("ifconfig", "eth0", ip, "netmask", netmask.NetMaskString(mask))
		}
		x("route", "add", "default", "gw", gw)
		resolv, err := getResolv()
		if err != nil {
			log.Fatal(err)
		}
		defer resolv.Close()
		err = writeDNS(resolv, dns)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func x(cmd string, args ...string) {
	if os.Getenv("USER") == "root" {
		cmd := exec.Command(cmd, args...)
		cmd.Stdout = os.Stdout
		cmd.Run()
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

func getDNSServers() []string {
	var servers []string
	for {
		ti := textinput.New("DNS server")
		ti.Validate = nil
		str, err := ti.RunPrompt()
		if err != nil {
			log.Fatal(err)
		}
		if str == "" {
			break
		}
		servers = append(servers, str)
	}
	return servers
}

func getResolv() (io.WriteCloser, error) {
	if os.Getenv("USER") != "root" {
		log.Info("not root, so writing to stdout")
		return os.Stdout, nil
	}
	resolv, err := os.Create("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("opening resolv.conf: %w", err)
	}
	return resolv, nil
}

func writeDNS(w io.Writer, servers []string) error {
	for _, srv := range servers {
		_, err := fmt.Fprintf(w, "nameserver %s\n", srv)
		if err != nil {
			return fmt.Errorf("writing dns server: %w", err)
		}
	}
	return nil
}
