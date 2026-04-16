package service

import accountv1 "github.com/solardome/gamepulse-platform/gen/account/v1"

type AccountService struct{}

func New() *AccountService {
	return &AccountService{}
}

func (s *AccountService) Ping(origin string) *accountv1.PingResponse {
	return &accountv1.PingResponse{
		Message: "OK",
		Service: "accounts-service",
	}
}
