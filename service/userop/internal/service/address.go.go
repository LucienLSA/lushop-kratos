package service

import (
	"context"
	v1 "userop/api/userop/v1"
	"userop/internal/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserOpService) GetAddressList(ctx context.Context, req *v1.AddressRequest) (*v1.AddressListResponse, error) {
	filter := domain.Address{
		UserID:       req.UserId,
		ID:           req.Id,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
	}
	listResp, err := s.as.GetAddressList(ctx, filter)
	if err != nil {
		return nil, err
	}
	out := &v1.AddressListResponse{Total: int32(listResp.Total)}
	for _, item := range listResp.List {
		out.Data = append(out.Data, &v1.AddressResponse{
			UserId:       item.UserID,
			Id:           item.ID,
			Province:     item.Province,
			City:         item.City,
			District:     item.District,
			Address:      item.Address,
			SignerName:   item.SignerName,
			SignerMobile: item.SignerMobile,
		})
	}
	return out, nil
}

func (s *UserOpService) CreateAddress(ctx context.Context, req *v1.AddressRequest) (*v1.AddressResponse, error) {
    addr := domain.Address{
        UserID:       req.UserId,
        ID:           req.Id,
        Province:     req.Province,
        City:         req.City,
        District:     req.District,
        Address:      req.Address,
        SignerName:   req.SignerName,
        SignerMobile: req.SignerMobile,
    }
    if err := s.as.CreateAddress(ctx, addr); err != nil {
        return nil, err
    }
    // Note: repo returns only error, so ID may not be populated here.
    return &v1.AddressResponse{
        UserId:       addr.UserID,
        Id:           addr.ID,
        Province:     addr.Province,
        City:         addr.City,
        District:     addr.District,
        Address:      addr.Address,
        SignerName:   addr.SignerName,
        SignerMobile: addr.SignerMobile,
    }, nil
}

func (s *UserOpService) DeleteAddress(ctx context.Context, req *v1.AddressRequest) (*emptypb.Empty, error) {
    addr := domain.Address{
        ID:           req.Id,
        UserID:       req.UserId,
        Province:     req.Province,
        City:         req.City,
        District:     req.District,
        Address:      req.Address,
        SignerName:   req.SignerName,
        SignerMobile: req.SignerMobile,
    }
    if err := s.as.DeleteAddress(ctx, addr); err != nil {
        return nil, err
    }
    return &emptypb.Empty{}, nil
}

func (s *UserOpService) UpdateAddress(ctx context.Context, req *v1.AddressRequest) (*emptypb.Empty, error) {
    addr := domain.Address{
        ID:           req.Id,
        UserID:       req.UserId,
        Province:     req.Province,
        City:         req.City,
        District:     req.District,
        Address:      req.Address,
        SignerName:   req.SignerName,
        SignerMobile: req.SignerMobile,
    }
    if err := s.as.UpdateAddress(ctx, addr); err != nil {
        return nil, err
    }
    return &emptypb.Empty{}, nil
}

func (s *UserOpService) MessageList(ctx context.Context, req *v1.MessageRequest) (*v1.MessageListResponse, error) {
    filter := domain.Message{
        ID:          req.Id,
        UserID:      req.UserId,
        MessageType: req.MessageType,
        Subject:     req.Subject,
        Message:     req.Message,
        File:        req.File,
    }
    listResp, err := s.as.GetMessageList(ctx, filter)
    if err != nil {
        return nil, err
    }
    out := &v1.MessageListResponse{Total: int32(listResp.Total)}
    for _, item := range listResp.List {
        out.Data = append(out.Data, &v1.MessageResponse{
            Id:          item.ID,
            UserId:      item.UserID,
            MessageType: item.MessageType,
            Subject:     item.Subject,
            Message:     item.Message,
            File:        item.File,
        })
    }
    return out, nil
}

func (s *UserOpService) CreateMessage(ctx context.Context, req *v1.MessageRequest) (*v1.MessageResponse, error) {
    msg := domain.Message{
        UserID:      req.UserId,
        MessageType: req.MessageType,
        Subject:     req.Subject,
        Message:     req.Message,
        File:        req.File,
    }
    if err := s.as.CreateMessage(ctx, msg); err != nil {
        return nil, err
    }
    return &v1.MessageResponse{
        Id:          msg.ID,
        UserId:      msg.UserID,
        MessageType: msg.MessageType,
        Subject:     msg.Subject,
        Message:     msg.Message,
        File:        msg.File,
    }, nil
}
