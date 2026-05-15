package main

import (
	"AppFir/controler"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	controler.Router()
}
