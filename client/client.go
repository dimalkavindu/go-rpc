// client implements a RPC client that can connect
// to a RPC server by either using plain TCP, HTTP
// or JSON-RPC.
//
// It consumes the `core` package to communicate
// with the server.
package client

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"strconv"

	"github.com/dimalkavindu/go-rpc/core"
	"github.com/dimalkavindu/go-rpc/menu"
)

// Client contains the configuration options for
// a RPC client that communicates with a RPC server
// over the network.
//
// Its parameters should match the server, for instance,
// if the server is offered via HTTP, it should have
// the property UseHttp set to true.
type Client struct {
	Port    uint
	UseHttp bool
	UseJson bool
	client  *rpc.Client
}

// Init initializes the underlying RPC client that is
// responsible for taking a codec and writing the RPC
// details down to it.
//
// Here we're using basic `(json)rpc.Dial*` but we could
// instead make use of raw TCP connections and tune them
// accordingly if we'd take this to production as we could
// then specify various timeouts and options for the transport
// layer.
//
// Note.: the HTTP thing is just a very thin layer of HTTP
// that is sent via the TCP connection. We could ditch that
// of and replace by a `CONNECT` call followed by checking
// the HTTP response that we got back.
//
// Note.: we're not setting TLS here either but it's a very
// simple thing given that we can have total control over
// the underlying connection.
var client *Client

func (c *Client) Init() (err error) {
	client = c
	if client.Port == 0 {
		err = errors.New("client: port must be specified")
		return
	}

	addr := "127.0.0.1:" + strconv.Itoa(int(client.Port))

	if client.UseHttp {
		client.client, err = rpc.DialHTTP("tcp", addr)
	} else if client.UseJson {
		client.client, err = jsonrpc.Dial("tcp", addr)
	} else {
		client.client, err = rpc.Dial("tcp", addr)
	}
	if err != nil {
		return
	}

	return
}

// Close gracefully terminates the underlying client.
func (c *Client) Close() (err error) {
	if client.client != nil {
		err = client.client.Close()
		return
	}

	return
}

func showVegitable(args ...string) error {
	var (
		request  = &core.Request{Command: args}
		response = new(core.Response)
	)

	err := client.client.Call("Handler.CshowVegitable", request, response)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if !response.Ok {
		fmt.Println(response.Message)
		return nil
	}

	if args[0] == "vegitable" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Vegitable Name", "Unit Price", "Stocks(KG)"})

		if args[1] == "all" {
			for _, v := range response.Vegitables.Vegitables {
				table.Append([]string{v.Name, v.PricePerKg, v.RemainingKgs})
			}
			table.Render()
			return nil
		} else {
			var v = response.Vegitables.Vegitables[0]
			table.Append([]string{v.Name, v.PricePerKg, v.RemainingKgs})
			table.Render()
			return nil
		}
	} else if args[0] == "price" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Vegitable Name", "Unit Price"})

		for _, v := range response.Vegitables.Vegitables {
			if v.Name == args[1] {
				table.Append([]string{v.Name, v.PricePerKg})
			}
		}
		table.Render()
		return nil
	} else if args[0] == "stocks" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Vegitable Name", "Stocks(KG)"})

		for _, v := range response.Vegitables.Vegitables {
			if v.Name == args[1] {
				table.Append([]string{v.Name, v.RemainingKgs})
			}
		}
		table.Render()
		return nil
	}

	return nil
}

func addVegitable(args ...string) error {
	var (
		request  = &core.Request{Command: args}
		response = new(core.Response)
	)

	err := client.client.Call("Handler.CaddVegitable", request, response)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	fmt.Println(response.Message)
	return nil
}

func updateVegitable(args ...string) error {
	var (
		request  = &core.Request{Command: args}
		response = new(core.Response)
	)

	err := client.client.Call("Handler.CupdateVegitable", request, response)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	fmt.Println(response.Message)
	return nil
}

func (c *Client) Start() (err error) {
	commandOptions := []menu.CommandOption{
		menu.CommandOption{"show", "\n" +
			"\tshow vegitable all\t: Shows all the vegitables\n" +
			"\tshow vegitable <vegitable name>\t: Shows unit price and stocks of a given vegitable\n" +
			"\tshow price <vegitable name>\t: Shows the unit price of a given vegitable\n" +
			"\tshow stocks <vegitable name>\t: Shows the stocks of a given vegitable", showVegitable},
		menu.CommandOption{"add", "\n" +
			"\tadd vegitable <vegitable name> <unit price> <stocks(KG)>\t: Adds a new vegitable with a given unit price & a stock value in KG", addVegitable},
		menu.CommandOption{"update", "\n" +
			"\tupdate price <vegitable name> <unit price>\t: Updates the unit price of a given vegitable\n" +
			"\tupdate stocks <vegitable name> <stocks(KG)>\t: Updates the stocks of a given vegitable", updateVegitable},
	}

	menuOptions := menu.NewMenuOptions("'menu' for help > ", 500)

	menu := menu.NewMenu(commandOptions, menuOptions)
	menu.Start()

	return
}
