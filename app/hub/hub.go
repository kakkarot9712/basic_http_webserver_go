package hub

import (
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
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
					if req.Method == "GET" {
						fileName := strings.TrimPrefix(req.Path, "/files/")
						buff, err := os.ReadFile(fmt.Sprintf("%s/%s", staticDirPath, fileName))
						if err != nil {
							ctx.Write(404, nil)
						}
						ctx.Headers().Set("Content-Type", "application/octet-stream")
						ctx.Headers().Set("Content-Length", fmt.Sprintf("%v", len(buff)))
						ctx.Write(200, buff)
					} else if req.Method == "POST" {
						fileName := strings.TrimPrefix(req.Path, "/files/")
						if req.Headers.Get("Content-Type") == "application/octet-stream" {
							bodySizeStr := req.Headers.Get("Content-Length")
							if bodySize, err := strconv.ParseUint(bodySizeStr, 10, 64); err == nil {
								if len(req.Body) != int(bodySize) {
									ctx.Write(400, nil)
								}
								file, err := os.Create(fmt.Sprintf("%s/%s", staticDirPath, fileName))
								if err != nil {
									ctx.Write(500, nil)
								}
								file.Write(req.Body)
								file.Sync() // Propagate changes to the disk from page cache
								ctx.Write(201, nil)
							} else {
								fmt.Println(err)
								ctx.Write(400, nil)
							}
						}
					} else {
						ctx.Write(404, nil)
					}

				default:
					ctx.Write(404, nil)
				}
			}
		}()
	}
	h.Wg = &wg
}
