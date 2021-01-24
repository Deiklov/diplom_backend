package api

import "github.com/Deiklov/diplom_backend/internal/services/api/server"

func main() {
	server := server.NewServer("localhost", 8080)
	server.Run()
}
