package main

import (
	"log"

	"github.com/colt3k/nglog/ng"
)

var escapeSequence = "\x1b"

func main() {
	log.Println("\x1b[30mtext\x1b[0m")
}

func bigTest() {
	someval := 1
	log.Println(ng.Black("some black text"))
	log.Println(ng.HiBlack("some hi black text"))
	log.Println(ng.Red("some red text %d", someval))
	log.Println(ng.Green("some green text"))
	log.Println(ng.Yellow("some yellow text"))
	log.Println(ng.Magenta("some magenta text"))
	log.Println(ng.Cyan("some cyan text"))
	log.Println(ng.White("some white text"))
	log.Println(ng.New(ng.FgBlack).Printf("black normal"))
	log.Println(ng.New(ng.FgBlack, ng.ITALIC).Printf("black italic"))
	log.Println(ng.New(ng.FgBlack, ng.BOLD).Printf("black bold"))
	log.Println(ng.New(ng.FgBlack, ng.UNDERLINE).Printf("black underline"))
	log.Println(ng.New(ng.FgBlack, ng.Reversed).Printf("black reversed"))
}
