// server Implements an RPC server that exposes
// a hello-world to RPC clients.
//
// The server is accessible either via HTTP, TCP
// or JSON-RPC.
package server

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"strconv"
	"time"

	"github.com/dimalkavindu/go-rpc/core"
	"github.com/dimalkavindu/go-rpc/menu"
	"github.com/olekukonko/tablewriter"
)

// Server holds the configuration used to initiate
// an RPC server.
type Server struct {
	Port     uint
	UseHttp  bool
	UseJson  bool
	Sleep    time.Duration
	listener net.Listener
}

// Close gracefully terminates the server listener.
func (s *Server) Close() (err error) {
	if s.listener != nil {
		err = s.listener.Close()
	}

	return
}

// initializing our vegitables array
var vegitables core.Vegitables

//Server busy flag
var serverBusy = false

func showVegitable(args ...string) error {
	if args[0] == "vegitable" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Vegitable Name", "Unit Price", "Stocks(KG)"})

		if args[1] == "all" {
			for _, v := range vegitables.Vegitables {
				table.Append([]string{v.Name, v.PricePerKg, v.RemainingKgs})
			}
			table.Render()
			return nil
		} else {
			for _, v := range vegitables.Vegitables {
				if v.Name == args[1] {
					table.Append([]string{v.Name, v.PricePerKg, v.RemainingKgs})
					table.Render()
					return nil
				}
			}
			fmt.Println("Vegitable '" + args[1] + "' is not found!")
			return nil
		}
	} else if args[0] == "price" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Vegitable Name", "Unit Price"})

		for _, v := range vegitables.Vegitables {
			if v.Name == args[1] {
				table.Append([]string{v.Name, v.PricePerKg})
				table.Render()
				return nil
			}
		}
		fmt.Println("Vegitable '" + args[1] + "' is not found!")
		return nil
	} else if args[1] == "stocks" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Vegitable Name", "Stocks(KG)"})

		for _, v := range vegitables.Vegitables {
			if v.Name == args[1] {
				table.Append([]string{v.Name, v.RemainingKgs})
				table.Render()
				return nil
			}
		}
		fmt.Println("Vegitable '" + args[1] + "' is not found!")
		return nil
	} else {
		fmt.Println("Unknown command format: " + args[0] + " " + args[1])
		return nil
	}
	return nil
}

func addVegitable(args ...string) error {
	if args[0] == "vegitable" {
		if len(args) == 4 {
			for _, v := range vegitables.Vegitables {
				if v.Name == args[1] {
					fmt.Println("Cannot add the vegitable! Vegitable '" + v.Name + "' already exists!")
					return nil
				}
			}

			var vegitable core.Vegitable
			vegitable.Name = args[1]
			vegitable.PricePerKg = args[2]
			vegitable.RemainingKgs = args[3]

			if serverBusy {
				fmt.Println("Server is busy! Please try again later!")
				return nil
			} else {
				serverBusy = true
			}

			vegitables.Vegitables = append(vegitables.Vegitables, vegitable)
			writeToDB()
			serverBusy = false
			fmt.Println("Vegitable '" + vegitable.Name + "' is successfully added")
			return nil
		} else {
			fmt.Println("Invalid number of inputs for 'add vegitable' command!")
			return nil
		}
	} else {
		fmt.Println("Unknown command format: 'add " + args[0] + "'")
		return nil
	}
	return nil
}

func updateVegitable(args ...string) error {
	if args[0] == "price" {
		if len(args) == 3 {
			for v := 0; v < len(vegitables.Vegitables); v++ {
				if vegitables.Vegitables[v].Name == args[1] {
					if serverBusy {
						fmt.Println("Server is busy! Please try again later!")
						return nil
					} else {
						serverBusy = true
					}

					vegitables.Vegitables[v].PricePerKg = args[2]
					writeToDB()
					serverBusy = false
					fmt.Println("Vegitable '" + vegitables.Vegitables[v].Name + "' is successfully updated!")
					return nil
				}
			}
			fmt.Println("Vegitable '" + args[1] + "' is not found!")
			return nil
		} else {
			fmt.Println("Invalid number of inputs for 'update price' command!")
			return nil
		}
	} else if args[0] == "stocks" {
		if len(args) == 3 {
			for v := 0; v < len(vegitables.Vegitables); v++ {
				if vegitables.Vegitables[v].Name == args[1] {
					if serverBusy {
						fmt.Println("Server is busy! Please try again later!")
						return nil
					} else {
						serverBusy = true
					}

					vegitables.Vegitables[v].RemainingKgs = args[2]
					writeToDB()
					serverBusy = false
					fmt.Println("Vegitable '" + vegitables.Vegitables[v].Name + "' is successfully updated!")
					return nil
				}
			}
			fmt.Println("Vegitable '" + args[1] + "' is not found!")
			return nil
		} else {
			fmt.Println("Invalid number of inputs for 'update stocks' command!")
			return nil
		}
	} else {
		fmt.Println("Unknown command format: 'add " + args[0] + "'")
		return nil
	}

	return nil
}

func writeToDB() error {
	// Open our xmlFile
	xmlFile, err := os.OpenFile("db.xml", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}

	encoder := xml.NewEncoder(xmlFile)
	err = encoder.Encode(vegitables)
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	return nil
}

// Handler holds the methods to be exposed by the RPC
// server as well as properties that modify the methods'
// behavior.
type Handler struct {
	// Sleep adds a little sleep between to the
	// method execution to simulate a time-consuming
	// operation.
	Sleep time.Duration
}

func (h *Handler) CshowVegitable(req core.Request, res *core.Response) (err error) {
	if req.Command[0] == "" {
		err = errors.New("Command shoud be specified!")
		return
	}

	if req.Command[0] == "vegitable" {
		if req.Command[1] == "all" {
			res.Message = "Command executed successfully!"
			res.Ok = true
			res.Vegitables = vegitables
		} else {
			for _, v := range vegitables.Vegitables {
				if v.Name == req.Command[1] {
					res.Message = "Command executed successfully!"
					res.Ok = true
					res.Vegitables.Vegitables = append(res.Vegitables.Vegitables, v)
					return nil
				}
			}
			res.Message = "Vegitable '" + req.Command[1] + "' is not found!"
			res.Ok = false
			return nil
		}
	} else if req.Command[0] == "price" || req.Command[0] == "stocks" {
		for _, v := range vegitables.Vegitables {
			if v.Name == req.Command[1] {
				res.Message = "Command executed successfully!"
				res.Ok = true
				res.Vegitables.Vegitables = append(res.Vegitables.Vegitables, v)
				return nil
			}
		}
		res.Message = "Vegitable '" + req.Command[1] + "' is not found!"
		res.Ok = false
		return nil
	} else {
		res.Message = "Unknown command format: " + req.Command[0] + " " + req.Command[1]
		res.Ok = false
		return nil
	}
	if h.Sleep != 0 {
		time.Sleep(h.Sleep)
	}
	return
}

func (h *Handler) CaddVegitable(req core.Request, res *core.Response) (err error) {
	if req.Command[0] == "" {
		err = errors.New("Command shoud be specified!")
		return
	}

	if req.Command[0] == "vegitable" {
		if len(req.Command) == 4 {
			for _, v := range vegitables.Vegitables {
				if v.Name == req.Command[1] {
					res.Message = "Cannot add the vegitable! Vegitable '" + v.Name + "' already exists!"
					res.Ok = false
					return nil
				}
			}

			var vegitable core.Vegitable
			vegitable.Name = req.Command[1]
			vegitable.PricePerKg = req.Command[2]
			vegitable.RemainingKgs = req.Command[3]
			if serverBusy {
				res.Message = "Server is busy! Please try again later!"
				res.Ok = false
				return nil
			} else {
				serverBusy = true
			}

			vegitables.Vegitables = append(vegitables.Vegitables, vegitable)
			writeToDB()
			serverBusy = false
			res.Message = "Vegitable '" + vegitable.Name + "' is added successfully!"
			res.Ok = true
			return nil
		} else {
			res.Message = "Invalid number of inputs for 'add vegitable' command!"
			res.Ok = false
			return nil
		}
	} else {
		res.Message = "Unknown command format: 'add " + req.Command[0] + "'"
		res.Ok = false
		return nil
	}
	return nil
}

func (h *Handler) CupdateVegitable(req core.Request, res *core.Response) (err error) {
	if req.Command[0] == "" {
		err = errors.New("Command shoud be specified!")
		return
	}

	if req.Command[0] == "price" {
		if len(req.Command) == 3 {
			for v := 0; v < len(vegitables.Vegitables); v++ {
				if vegitables.Vegitables[v].Name == req.Command[1] {
					if serverBusy {
						res.Message = "Server is busy! Please try again later!"
						res.Ok = false
						return nil
					} else {
						serverBusy = true
					}

					vegitables.Vegitables[v].PricePerKg = req.Command[2]
					writeToDB()
					serverBusy = false
					res.Message = "Vegitable '" + vegitables.Vegitables[v].Name + "' is updated successfully!"
					res.Ok = true
					return nil
				}
			}
			res.Message = "Vegitable '" + req.Command[1] + "' is not found!"
			res.Ok = false
			return nil
		} else {
			res.Message = "Invalid number of inputs for 'update price' command!"
			res.Ok = false
			return nil
		}
	} else if req.Command[0] == "stocks" {
		if len(req.Command) == 3 {
			for v := 0; v < len(vegitables.Vegitables); v++ {
				if vegitables.Vegitables[v].Name == req.Command[1] {
					if serverBusy {
						res.Message = "Server is busy! Please try again later!"
						res.Ok = false
						return nil
					} else {
						serverBusy = true
					}

					vegitables.Vegitables[v].RemainingKgs = req.Command[2]
					writeToDB()
					serverBusy = false
					res.Message = "Vegitable '" + vegitables.Vegitables[v].Name + "' is updated successfully!"
					res.Ok = true
					return nil
				}
			}
			res.Message = "Vegitable '" + req.Command[1] + "' is not found!"
			res.Ok = false
			return nil
		} else {
			res.Message = "Invalid number of inputs for 'update price' command!"
			res.Ok = false
			return nil
		}
	} else {
		res.Message = "Unknown command format: 'add " + req.Command[0] + "'"
		res.Ok = false
		return nil
	}
	return nil
}

// Starts initializes the RPC server by first verifying
// if all the necessary configuration has been set.
//
// It then publishes the receiver's methods (core.Handler)
// in the default RPC server. By doing so, the `Handler`
// public methods that satisfy the rpc interface become
// available to clients connecting to this server.
//
// With the receiver registered, it starts the server
// such that new connections can be accepted.
func (s *Server) StartServer() (err error) {
	if s.Port <= 0 {
		err = errors.New("port must be specified")
		return
	}

	rpc.Register(&Handler{
		Sleep: s.Sleep,
	})

	s.listener, err = net.Listen("tcp", ":"+strconv.Itoa(int(s.Port)))
	if err != nil {
		return
	}

	// Open our xmlFile
	xmlFile, err := os.Open("db.xml")
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// unmarshaling  byteArray
	xml.Unmarshal(byteValue, &vegitables)

	go func() {
		s.StartMenu()
		s.Close()
		os.Exit(0)
	}()

	if s.UseHttp {
		rpc.HandleHTTP()
		http.Serve(s.listener, nil)
	} else if s.UseJson {
		var conn net.Conn

		for {
			conn, err = s.listener.Accept()
			if err != nil {
				return
			}

			jsonrpc.ServeConn(conn)
		}

	} else {
		rpc.Accept(s.listener)
	}

	return
}

func (s *Server) StartMenu() (err error) {
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
