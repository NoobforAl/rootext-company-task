package main

import "ratblog/cmd"

func main() {
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
