package main

import (
	"github.com/go-kratos/kratos/v2/log"
	kitgen "github.com/liujitcn/gorm-kit/gen"
	_ "github.com/liujitcn/kratos-kit/database/gorm/driver/mysql"
)

const defaultDSN = "root:112233@tcp(127.0.0.1:3306)/shop?charset=utf8&parseTime=True&loc=Local&timeout=1000ms"

func main() {
	g := kitgen.NewGen(
		kitgen.WithDriver("mysql"),
		kitgen.WithSource(defaultDSN),
	)
	if err := g.Execute(); err != nil {
		log.Fatal(err)
	}
}
