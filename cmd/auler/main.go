package main

import (
	"os"

	// 导入匿名宝automaxprocs自动配额CPU
	_ "go.uber.org/automaxprocs"

	"github.com/HeapSoil/auler/internal/auler"
)

func main() {
	command := auler.NewAulerCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}

}
