package main

import (
    "WedhaWS/utsgosmt5/alumni-tracer/injector"
    "WedhaWS/utsgosmt5/alumni-tracer/routes"
    "github.com/gofiber/fiber/v2/log"
    "os"
    "io"
) 	

func main() {
    app := injector.InitializedService()
    router := routes.NewRouter(app)
    file, _ := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    iw := io.MultiWriter(os.Stdout, file)
    log.SetOutput(iw)
    log.Fatal(router.Listen(":3000"))
}
