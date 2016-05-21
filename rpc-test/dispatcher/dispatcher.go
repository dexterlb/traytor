package dispatcher

import (
	"log"

	"github.com/valyala/gorpc"
)

type NumberList struct {
	Data []int
}

func Sum(input *NumberList) (int, error) {
	log.Printf("adding %d numbers...", len(input.Data))
	reply := 0
	for i := range input.Data {
		reply = reply + input.Data[i]
	}
	log.Printf("done: %d", reply)
	return reply, nil
}

func Register() {
	gorpc.RegisterType(&NumberList{})
}

func CreateDispatcher() *gorpc.Dispatcher {
	d := gorpc.NewDispatcher()
	d.AddFunc("Sum", Sum)
	Register()
	return d
}
