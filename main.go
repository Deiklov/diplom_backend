package main

import "github.com/Deiklov/diplom_backend/internal/services/api/server"

func main() {
	serverObj := server.NewServer("localhost", 8080)
	serverObj.Run()
}
