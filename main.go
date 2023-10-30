package main

import (
    "web-filter/routes"
    "web-filter/database"
)

// @title WebFilter - API
// @version 1.0
// @description Gerencia e lida com as escritas para o serviço Squid - WebFilter
// @host 172.23.58.10:8080
// @BasePath /webfilter
// @schemes http https
func main() {
    database.ConectaNoBD()

    routes.HandleRequests()
}