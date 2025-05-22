package hub

import (
	"net"
	"strings"
	"sync"

	"github.com/codecrafters-io/http-server-starter-go/app/gcho"
)

type Task struct {
	C net.Conn
}

type Hub struct {
	workersLimit uint
	workChann    chan *Task
	WorkReceiver chan<- *Task
	Wg           *sync.WaitGroup
}

func NewHub(wsize uint) *Hub {
	workChann := make(chan *Task, wsize)
	return &Hub{
		workersLimit: wsize,
		workChann:    workChann,
		WorkReceiver: workChann,
	}
}

func (h *Hub) Start() {
	// create workers
	var wg sync.WaitGroup

	for range h.workersLimit {
		wg.Add(1)
		go func() {
			for task := range h.workChann {
				ctx, err := gcho.NewContext(task.C)
				if err != nil {
					panic(err)
				}
				req := ctx.Request
				switch {
				case req.Path == "/":
					ctx.Write(200, nil)
				case strings.HasPrefix(req.Path, "/echo"):
					echoStr := strings.TrimPrefix(req.Path, "/echo/")
					ctx.Write(200, []byte(echoStr))
				case req.Path == "/user-agent":
					ctx.Write(200, []byte(ctx.Request.Headers.Get("User-Agent")))
				default:
					ctx.Write(404, nil)
				}
			}
		}()
	}
	h.Wg = &wg
}
