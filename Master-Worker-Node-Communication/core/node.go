package core

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type MasterNode struct {
	api     *gin.Engine
	ln      net.Listener
	svr     *grpc.Server
	nodeSvr *NodeServiceGrpcServer
}

func (n *MasterNode) Init() (err error) {
	n.ln, err = net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	n.svr = grpc.NewServer()
	n.nodeSvr = GetNodeServiceGrpcServer()
	RegisterNodeServiceServer(node.svr, n.nodeSvr)

	n.api = gin.Default()
	n.api.POST("/tasks", func(c *gin.Context) {
		var payload struct {
			Cmd string `json:"cmd"`
		}
		if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		n.nodeSvr.CmdChannel <- payload.Cmd
		c.AbortWithStatusJSON(200, http.StatusOK)
	})
	return nil
}

func (n *MasterNode) Start() {
	go n.svr.Serve(n.ln)
	_ = n.api.Run(":9092")
	n.svr.Stop()
}

var node *MasterNode

func GetMasterNode() *MasterNode {
	if node == nil {
		node = &MasterNode{}
		if err := node.Init(); err != nil {
			panic(err)
		}
	}
	return node
}
