package main

import "redirector/app"

func main() {
	if err := app.New().Run(); err != nil {
		_ = err
	}
}
