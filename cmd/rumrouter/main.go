package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/rumis/rumrouter-go/pkg/config"
	"github.com/rumis/rumrouter-go/pkg/generator"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("current work directory is not available:", err)
		return
	}
	var frameType = flag.String("t", "echov4", "framework type, now echo,echov4,gin is available")
	var packName = flag.String("p", "rumrouter", "package name,default is rumrouter")
	var inputSourcePath = flag.String("i", dir, "input source code dirctory,default is current work dirctory")
	var outputPath = flag.String("o", path.Join(dir, "rumrouter"), "router code dirctory,default is rumrouter")
	var namespace = flag.String("n", "", "namespace, default is empty ,only namespace in defined and matched,the router will be init")
	flag.Parse()

	opts := config.Options{
		FrameworkType: *frameType,
		PackageName:   *packName,
		SourcePath:    *inputSourcePath,
		OutputPath:    *outputPath,
		Namespace:     *namespace,
	}

	err = generator.InitRouter(opts)
	if err != nil {
		fmt.Println("router initialize failed:", err)
		return
	}
	fmt.Println("router initialize success")
}
