package main

type test struct {
	ip   string
	port int
}

func main() {
	t := map[int]test{1: {ip: "", port: 0}}
	t[2] = test{ip: "", port: 0}
}
