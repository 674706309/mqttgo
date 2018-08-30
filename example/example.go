package main

import "goim/mqtt"

func main() {
	conn := mqtt.NewConnack()
	conn.GetType()
}
