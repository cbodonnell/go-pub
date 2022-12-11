package workers

type Worker interface {
	Start()
	GetChannel() chan interface{}
}
