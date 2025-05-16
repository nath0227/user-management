package user

import (
	"context"
	usergrpc "user-management/app/user/grpc/gen/go/user/v1"
	"user-management/response"
)

type GrpcHandler struct {
	usecase Usecase
}

func NewGrpcHandler(u Usecase) *GrpcHandler {
	return &GrpcHandler{usecase: u}
}

func (h *GrpcHandler) CreateUser(ctx context.Context, req *usergrpc.CreateUserRequest) (*usergrpc.CreateUserResponse, error) {
	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.usecase.CreateUser(ctx, user)
	if err != nil {
		if err.Error() == EmailAlreadyExists {
			return &usergrpc.CreateUserResponse{
				Code:    response.DuplicatedRegistration().Code,
				Message: response.DuplicatedRegistration().Message,
			}, nil
		}
		return nil, err
	}

	return &usergrpc.CreateUserResponse{
		Code:    response.Success().Code,
		Message: response.Success().Message,
		Data: &usergrpc.CreateUserResponse_Data{
			Id: resp.Data.(CreateResponse).Id,
		},
	}, nil
}

func (h *GrpcHandler) GetUser(ctx context.Context, req *usergrpc.GetUserRequest) (*usergrpc.GetUserResponse, error) {

	if respValidate := IdValidation(req.Id); !respValidate.IsSuccess() {
		resp := response.InvalidData("id")
		return &usergrpc.GetUserResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}

	resp, err := h.usecase.FindUserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return &usergrpc.GetUserResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}

	user := resp.Data.(User)
	return &usergrpc.GetUserResponse{
		Code:    response.Success().Code,
		Message: response.Success().Message,
		Data: &usergrpc.GetUserResponse_Data{
			Id:    user.ID.Hex(),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
