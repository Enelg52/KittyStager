package main

import (
	s "github.com/inancgumus/prettyslice"
)

func main() {
	nums := string("sdf")

	// Render to stdout by default
	s.Show("nums", nums)
}
