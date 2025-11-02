// grpc endpoint handler.
package grpcserver

import (
	"context"
	"evolve/db/connection"
	"evolve/proto"
	"evolve/util/auth"
	dbutil "evolve/util/db/user"
	"os"
)

type GRPCServer struct {
	proto.UnimplementedAuthenticateServer
}

func (*GRPCServer) Auth(ctx context.Context, req *proto.TokenValidateRequest) (*proto.TokenValidateResponse, error) {
	user, err := auth.ValidateToken(req.GetToken())

	if err != nil {
		return &proto.TokenValidateResponse{
			Valid: false,
		}, nil
	}

	db, err := connection.PoolConn(ctx)
	if err != nil {
		return &proto.TokenValidateResponse{
			Valid: false,
		}, err
	}

	//checking for CSRF token
	if os.Getenv("CSRF_PROTECTION") == "true" && req.GetCsrfToken() != "" && user["csrf_token"] != "" {

		csrfToken := req.GetCsrfToken()
		if csrfToken != user["csrf_token"] {
			return &proto.TokenValidateResponse{
				Valid: false,
			}, nil
		}
	} else {
		return &proto.TokenValidateResponse{
			Valid: false,
		}, nil
	}

	userData, err := dbutil.UserById(ctx, user["id"], db)
	if err != nil {
		return &proto.TokenValidateResponse{
			Valid: false,
		}, err
	}

	return &proto.TokenValidateResponse{
		Valid:    true,
		Id:       userData["id"],
		UserName: userData["userName"],
		Email:    userData["email"],
		Role:     userData["role"],
		FullName: userData["fullName"],
	}, nil
}
