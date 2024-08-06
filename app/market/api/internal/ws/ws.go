package ws

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"net/http"
	"strings"
)

const ROOM = "market"

// WebSocketServer 把这个做成中间件
type WebSocketServer struct {
	prefix string
	server *socketio.Server
}

func (ws *WebSocketServer) Start() {
	if err := ws.server.Serve(); err != nil {
		logx.Error(err)
	}
}

func (ws *WebSocketServer) Stop() {
	if err := ws.server.Close(); err != nil {
		logx.Error(err)
	}
}

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func NewWebSocketServer(prefix string) *WebSocketServer {

	// 处理跨域
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	server.OnConnect("/", func(conn socketio.Conn) error {
		conn.SetContext("")
		log.Println("ws connected: ", conn.ID())
		conn.Join(ROOM)
		return nil
	})

	return &WebSocketServer{
		prefix: prefix,
		server: server,
	}
}

func (ws *WebSocketServer) ServerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		logx.Info("===================path: ", path)
		if strings.HasPrefix(path, ws.prefix) {
			// 说明是socket.io的请求
			ws.server.ServeHTTP(writer, request)
		} else {
			next.ServeHTTP(writer, request)
		}
	})
}

// BroadcastToRoom 向前端推送数据
// path: "/"
// event: "/topic/market/thumb"
func (ws *WebSocketServer) BroadcastToRoom(path string, event string, data any) {
	go func() {
		ws.server.BroadcastToRoom(path, ROOM, event, data)
		//log.Println("通过socket.io发送到前端的数据，event：", event)
		//log.Println("通过socket.io发送到前端的数据，data：", data.(string))
	}()
}
