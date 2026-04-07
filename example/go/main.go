package main

import (
	"context"
	"net/http"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"

	proto "github.com/sa3akash/proto-package-example/gen/go/users/v1"
	"github.com/sa3akash/proto-package-example/gen/go/users/v1/usersv1connect"
)

type UserServer struct{}

// DeleteUser implements [usersv1connect.UserServiceHandler].
func (s *UserServer) DeleteUser(context.Context, *connect.Request[proto.DeleteUserRequest]) (*connect.Response[proto.DeleteUserResponse], error) {
	panic("unimplemented")
}

// GetUser implements [usersv1connect.UserServiceHandler].
func (s *UserServer) GetUser(context.Context, *connect.Request[proto.GetUserRequest]) (*connect.Response[proto.GetUserResponse], error) {
	panic("unimplemented")
}

// ListUsers implements [usersv1connect.UserServiceHandler].
func (s *UserServer) ListUsers(context.Context, *connect.Request[proto.ListUsersRequest]) (*connect.Response[proto.ListUsersResponse], error) {
	panic("unimplemented")
}

// UpdateUser implements [usersv1connect.UserServiceHandler].
func (s *UserServer) UpdateUser(context.Context, *connect.Request[proto.UpdateUserRequest]) (*connect.Response[proto.UpdateUserResponse], error) {
	panic("unimplemented")
}



func (s *UserServer) CreateUser(
	ctx context.Context,
	req *connect.Request[proto.CreateUserRequest],
) (*connect.Response[proto.CreateUserResponse], error) {
	// Validate request fields against protovalidate constraints
	v, _ := protovalidate.New()
	if err := v.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	user := &proto.User{
		Id:    "generated-uuid",
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
		Role:  req.Msg.Role,
	}
	return connect.NewResponse(&proto.CreateUserResponse{User: user}), nil
}

func main() {
	mux := http.NewServeMux()
	path, handler := usersv1connect.NewUserServiceHandler(&UserServer{})
	mux.Handle(path, handler)
	http.ListenAndServe(":8080", mux)
}
