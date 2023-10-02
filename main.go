package main

import (
    "web-filter/routes"
    "web-filter/database"
)

func main() {
    database.ConectaNoBD()

    routes.HandleRequests()
}