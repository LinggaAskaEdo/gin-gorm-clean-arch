package main

import (
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/bootstrap"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	godotenv.Load()
	fx.New(bootstrap.Module).Run()
}
