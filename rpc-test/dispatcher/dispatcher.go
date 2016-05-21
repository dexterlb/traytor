package dispatcher

import (
	"github.com/valyala/gorpc"
)

type NumberList struct {
	Data []int
}

func Sum(input *NumberList) (int, error) {
	reply := 0
	for i := range input.Data {
		reply = reply + input.Data[i]
	}
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
