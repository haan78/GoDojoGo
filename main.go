package main

import (
	"GoDojoGo/deff"
	globals "GoDojoGo/deff"
	"GoDojoGo/libbrevo"
	"fmt"
	"os"
)

func main() {

	runWorker := true

	if err := globals.LoadSettings(".env"); err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "version":
			return
		case "noworker":
			runWorker = false
		case "brevocheck":
			rep, err := libbrevo.SendinblueCheck(deff.Settings.BREVO_API_KEY)
			if err == nil {
				for _, r := range rep {
					fmt.Println(r)
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
	}

	if runWorker {
		fmt.Println("Running with worker...")
		RunWithWorker()
	} else {
		fmt.Println("Running without worker...")
		RunWithoutWorker()
	}

}
