package server

type NumberList struct {
	Data []int
}

type Sumator int

func (s *Sumator) Sumatorirane(input *NumberList, reply *int) error {
	*reply = 0
	for i := range input.Data {
		*reply = *reply + input.Data[i]
	}
	return nil
}
