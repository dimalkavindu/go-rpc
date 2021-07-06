
        DESCRIPTION

                This is a simple distributed application for "Vegetable Sales System"
                as described below using "Go Language" and "Remote Procedure Calls (RPC)"

                Server will maintain a file which keeps records of different available vegetables
                including price per kg and available amount of kg for each vegetable.

                The server has following functions.

                    1. Query the file and output names of all available vegetables.
                    2. Output the price per kg of a given vegetable.
                    3. Output the available amount of kg for a given vegetable.
                    4. Add new vegetable to the file with price per kg and among of kg.
                    5. Update the price or available amount of a given vegetable.

                A client can use server functions to do the following tasks.

                    1. Receive a list of all available vegetables and display.
                    2. Get the price per kg of a given vegetable and display.
                    3. Get the available amount of kg of a given vegetable and display.
                    4. Send a new vegetable name to the server to be added to the server file.
                    5. Send new price or available amount for a given vegetable to be updated in the server file.

                Both client and server are meant to be run using a single binary
                    To run as a server, turn the `-server` flag on:

                        ./main.exe -server

                    To run as a client:

                        ./main.exe

                For the demonstration purposed, the server and client will be on the same node.


        USAGE

                ./main --help
                Usage of ./main:
                  -http
                        whether it should use HTTP
                  -json
                        whether it should use json-rpc
                  -port uint
                        port to listen or connect to for rpc calls (default 1337)
                  -server
                        activates server mode


