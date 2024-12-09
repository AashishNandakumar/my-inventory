package main

func main() {
	app := App{}
	app.Initialize("localhost", "5436", "admin", "admin", "inventory")
	app.Run("localhost:10000")
}
