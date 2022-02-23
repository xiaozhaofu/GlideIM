package api

import (
	"context"
	"go_im/im/api"
	"go_im/protobuff/pb_rpc"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	*rpc.BaseServer
}

func (s *Server) Handle(ctx context.Context, r *pb_rpc.ApiHandleRequest, resp *emptypb.Empty) error {

	api.Handle(r.GetUid(), r.GetDevice(), r.GetMessage())
	return nil
}

func (s *Server) Echo(ctx context.Context, r *pb_rpc.ApiHandleRequest, resp *pb.Response) error {
	return nil
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	s.Register(options.Name, s)
	return s
}
