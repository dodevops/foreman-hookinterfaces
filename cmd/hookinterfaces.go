package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/go-resty/resty/v2"
	"net"
	"os"
	"strings"
)

type Host struct {
	ID   int
	Name string
}
type HostResponse struct {
	Results []*Host
}

//goland:noinspection GoSnakeCaseUsage
type Interface struct {
	Subnet_ID  *int
	ID         int
	Identifier string
	Ip         string
}

type InterfaceResponse struct {
	Results []*Interface
}

type Todo struct {
	Host      Host
	Interface Interface
}

//goland:noinspection GoSnakeCaseUsage
type Subnet struct {
	Id              int
	Name            string
	Network_Address string
}

type SubnetResult struct {
	Results []*Subnet
}

func main() {
	parser := argparse.NewParser("hookinterfaces", "Hooks a Foreman host interface into the fitting subnet")
	baseUrl := parser.String("b", "base-url", &argparse.Options{
		Help:     "Base URL of foreman system to use",
		Required: true,
	})
	username := parser.String("u", "username", &argparse.Options{
		Help:     "Foreman admin user to use",
		Required: true,
	})
	password := parser.String("p", "password", &argparse.Options{
		Help:     "Foreman password to use",
		Required: true,
	})
	dryRun := parser.Flag("d", "dryrun", &argparse.Options{
		Help:    "Only print what would be done",
		Default: false,
	})
	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	hosts := &HostResponse{}
	client := resty.New().SetBasicAuth(*username, *password).SetBaseURL(*baseUrl)
	if _, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(hosts).
		SetQueryParams(map[string]string{"per_page": "all"}).
		Get("/api/hosts"); err != nil {
		panic(err)
	}

	var todos []Todo
	for _, host := range hosts.Results {
		interfaces := &InterfaceResponse{}
		if _, err := client.R().
			SetHeader("Accept", "application/json").
			SetResult(interfaces).
			Get(fmt.Sprintf("/api/hosts/%d/interfaces", host.ID)); err != nil {
			panic(err)
		}

		for _, hostInterface := range interfaces.Results {
			if hostInterface.Subnet_ID == nil {
				todos = append(todos, Todo{
					Host:      *host,
					Interface: *hostInterface,
				})
			}
		}
	}

	subnets := &SubnetResult{}
	if _, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(subnets).
		Get("/api/subnets"); err != nil {
		panic(err)
	}

TODO:
	for _, todo := range todos {
		if strings.HasPrefix(todo.Interface.Identifier, "br-") {
			continue TODO
		}
		if todo.Interface.Ip == "" {
			continue TODO
		}
		for _, subnet := range subnets.Results {
			if _, c, err := net.ParseCIDR(subnet.Network_Address); err != nil {
				panic(err)
			} else {
				if c.Contains(net.ParseIP(todo.Interface.Ip)) {
					println(fmt.Sprintf("Interface %s of host %s will be set to subnet %s", todo.Interface.Identifier, todo.Host.Name, subnet.Name))
					if !*dryRun {
						_, err := client.R().
							SetBody(map[string]map[string]int{
								"interface": {
									"subnet_id": subnet.Id,
								},
							}).
							Put(fmt.Sprintf("/api/hosts/%d/interfaces/%d", todo.Host.ID, todo.Interface.ID))
						if err != nil {
							panic(err)
						}
					}
					continue TODO
				}
			}
		}
		println(fmt.Sprintf("Did not find a subnet for interface %s of host %s with IP %s", todo.Interface.Identifier, todo.Host.Name, todo.Interface.Ip))
	}

}
