package hub

import (
	"fmt"
	"net"
	"os"
	"slices"
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
	dirArgIndex := slices.Index(os.Args, "--directory")
	var staticDirPath string
	if dirArgIndex != -1 {
		if len(os.Args) < dirArgIndex+1 {
			panic("directory location is required!")
		}
		staticDirPath = os.Args[dirArgIndex+1]
	}

	for range h.workersLimit {
		wg.Add(1)
		go func() {
			select {
			case task, ok := <-h.workChann:
				if !ok {
					wg.Done()
					return
				}
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
				case strings.HasPrefix(req.Path, "/files"):
					fileName := strings.TrimPrefix(req.Path, "/files/")
					buff, err := os.ReadFile(fmt.Sprintf("%s/%s", staticDirPath, fileName))
					if err != nil {
						ctx.Write(404, nil)
					}
					ctx.Headers().Set("Content-Type", "application/octet-stream")
					ctx.Headers().Set("Content-Length", fmt.Sprintf("%v", len(buff)))
					ctx.Write(200, buff)
				default:
					ctx.Write(404, nil)
				}
			}
		}()
	}
	h.Wg = &wg
}
