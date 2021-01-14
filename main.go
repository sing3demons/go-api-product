package main

import (
	"kp-app/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//สร้าง folder
	uploadDirs := [...]string{"products", "users"}
	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/"+dir, 0755)
	}

	routes.Serve(r)

	r.Run()
}
