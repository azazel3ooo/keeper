package main

import "github.com/azazel3ooo/keeper/internal/apps/server"

// @title           Swagger Keeper server
// @version         0.0.1
// @description     This is a image processing server.
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath        /
func main() {
	server.Start()
}
