// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: alert/alert.proto

package alert

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v3/api"
	server "github.com/micro/go-micro/v3/server"
	client "github.com/micro/micro/v3/service/client"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Alert service

func NewAlertEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Alert service

type AlertService interface {
	// ReportEvent does event ingestions.
	ReportEvent(ctx context.Context, in *ReportEventRequest, opts ...client.CallOption) (*ReportEventResponse, error)
}

type alertService struct {
	c    client.Client
	name string
}

func NewAlertService(name string, c client.Client) AlertService {
	return &alertService{
		c:    c,
		name: name,
	}
}

func (c *alertService) ReportEvent(ctx context.Context, in *ReportEventRequest, opts ...client.CallOption) (*ReportEventResponse, error) {
	req := c.c.NewRequest(c.name, "Alert.ReportEvent", in)
	out := new(ReportEventResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Alert service

type AlertHandler interface {
	// ReportEvent does event ingestions.
	ReportEvent(context.Context, *ReportEventRequest, *ReportEventResponse) error
}

func RegisterAlertHandler(s server.Server, hdlr AlertHandler, opts ...server.HandlerOption) error {
	type alert interface {
		ReportEvent(ctx context.Context, in *ReportEventRequest, out *ReportEventResponse) error
	}
	type Alert struct {
		alert
	}
	h := &alertHandler{hdlr}
	return s.Handle(s.NewHandler(&Alert{h}, opts...))
}

type alertHandler struct {
	AlertHandler
}

func (h *alertHandler) ReportEvent(ctx context.Context, in *ReportEventRequest, out *ReportEventResponse) error {
	return h.AlertHandler.ReportEvent(ctx, in, out)
}
