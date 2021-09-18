package route

import (
	"github.com/stretchr/testify/assert"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

var etcdSrv = []string{"127.0.0.1:2379", "127.0.0.1:2381", "127.0.0.1:2383"}

func TestNewServer(t *testing.T) {

	op := rpc.ServerOptions{
		Name:        "route",
		Network:     "tcp",
		Addr:        "127.0.0.1",
		Port:        8977,
		EtcdServers: etcdSrv,
	}

	routeServer := NewServer(&op)
	err := routeServer.Run()
	t.Error(err)
}

func TestNewServer2(t *testing.T) {

	op := rpc.ServerOptions{
		Name:        "route",
		Network:     "tcp",
		Addr:        "127.0.0.1",
		Port:        8976,
		EtcdServers: etcdSrv,
	}

	routeServer := NewServer(&op)
	err := routeServer.Run()
	t.Error(err)
}

func TestClient_Register(t *testing.T) {

	cli := newClient()
	defer cli.Close()
	err := cli.Register(&pb.RegisterRtReq{
		SrvId:           "api",
		SrvName:         "api",
		RoutePolicy:     1,
		DiscoverySrvUrl: etcdSrv,
		DiscoveryType:   1,
	}, &emptypb.Empty{})

	if err != nil {
		t.Error(err)
	}
}

func TestClient_Route(t *testing.T) {
	cli := newClient()
	defer cli.Close()

	req := &pb.HandleRequest{
		Uid: 1,
		Message: &pb.Message{
			Seq:    1,
			Action: "api.app.echo",
			Data:   "echo_test",
		},
	}
	resp := &pb.Response{}
	err := cli.RouteByTag("api.Echo", "", req, resp)
	assert.Nil(t, err)
	assert.Equal(t, req.Message.Data, resp.Message)
}

func TestServer_SetTag(t *testing.T) {
	cli := newClient()
	defer cli.Close()
	err := cli.SetTag("api", "uid_001", "tcp@127.0.0.1:8973")
	assert.Nil(t, err)
}

func TestClient_RemoveTag(t *testing.T) {
	cli := newClient()
	defer cli.Close()
	err := cli.RemoveTag("api", "uid_001")
	assert.Nil(t, err)
}

func newClient() *Client {
	client := NewClient(&rpc.ClientOptions{
		Name:        "route",
		EtcdServers: etcdSrv,
	})
	err := client.Run()
	if err != nil {
		panic(err)
	}
	return client
}