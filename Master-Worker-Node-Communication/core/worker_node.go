package core

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"google.golang.org/grpc"
)

type WokerNode struct {
	conn *grpc.ClientConn
	c    NodeServiceClient
}

func (n *WokerNode) Init() (err error) {
	n.conn, err = grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return err
	}

	n.c = NewNodeServiceClient(n.conn)

	return nil

}

func (n *WokerNode) Start() {
	fmt.Println("worker node started")

	_, _ = n.c.ReportStatus(context.Background(), &Request{})

	stream, _ := n.c.AssignTask(context.Background(), &Request{})
	for {
		res, err := stream.Recv()
		if err != nil {
			return
		}

		fmt.Print("received command: ", res.Data)

		parts := strings.Split(res.Data, " ")
		if err := exec.Command(parts[0], parts[1:]...).Run(); err != nil {
			fmt.Println(err)
		}
	}

}

var workerNode *WokerNode

func GetWorkerNode() *WokerNode {
	if workerNode == nil {
		workerNode = &WokerNode{}

		if err := workerNode.Init(); err != nil {
			panic(err)
		}
	}
	return workerNode
}
