package server

import (
	"context"

	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/router"
	pb "github.com/micro/micro/v2/service/router/proto"
)

type Table struct {
	Router router.Router
}

func (t *Table) Create(ctx context.Context, route *pb.Route, resp *pb.CreateResponse) error {
	err := t.Router.Table().Create(router.Route{
		Service:  route.Service,
		Address:  route.Address,
		Gateway:  route.Gateway,
		Network:  route.Network,
		Router:   route.Router,
		Link:     route.Link,
		Metric:   route.Metric,
		Metadata: route.Metadata,
	})
	if err != nil {
		return errors.InternalServerError("go.micro.router", "failed to create route: %s", err)
	}

	return nil
}

func (t *Table) Update(ctx context.Context, route *pb.Route, resp *pb.UpdateResponse) error {
	err := t.Router.Table().Update(router.Route{
		Service:  route.Service,
		Address:  route.Address,
		Gateway:  route.Gateway,
		Network:  route.Network,
		Router:   route.Router,
		Link:     route.Link,
		Metric:   route.Metric,
		Metadata: route.Metadata,
	})
	if err != nil {
		return errors.InternalServerError("go.micro.router", "failed to update route: %s", err)
	}

	return nil
}

func (t *Table) Delete(ctx context.Context, route *pb.Route, resp *pb.DeleteResponse) error {
	err := t.Router.Table().Delete(router.Route{
		Service:  route.Service,
		Address:  route.Address,
		Gateway:  route.Gateway,
		Network:  route.Network,
		Router:   route.Router,
		Link:     route.Link,
		Metric:   route.Metric,
		Metadata: route.Metadata,
	})
	if err != nil {
		return errors.InternalServerError("go.micro.router", "failed to delete route: %s", err)
	}

	return nil
}

// List returns all routes in the routing table
func (t *Table) List(ctx context.Context, req *pb.Request, resp *pb.ListResponse) error {
	routes, err := t.Router.Table().List()
	if err != nil {
		return errors.InternalServerError("go.micro.router", "failed to list routes: %s", err)
	}

	respRoutes := make([]*pb.Route, 0, len(routes))
	for _, route := range routes {
		respRoute := &pb.Route{
			Service:  route.Service,
			Address:  route.Address,
			Gateway:  route.Gateway,
			Network:  route.Network,
			Router:   route.Router,
			Link:     route.Link,
			Metric:   route.Metric,
			Metadata: route.Metadata,
		}
		respRoutes = append(respRoutes, respRoute)
	}

	resp.Routes = respRoutes

	return nil
}

func (t *Table) Query(ctx context.Context, req *pb.QueryRequest, resp *pb.QueryResponse) error {
	routes, err := t.Router.Table().Query(
		router.QueryService(req.Query.Service),
		router.QueryNetwork(req.Query.Network),
	)
	if err != nil {
		return errors.InternalServerError("go.micro.router", "failed to lookup routes: %s", err)
	}

	respRoutes := make([]*pb.Route, 0, len(routes))
	for _, route := range routes {
		respRoute := &pb.Route{
			Service:  route.Service,
			Address:  route.Address,
			Gateway:  route.Gateway,
			Network:  route.Network,
			Router:   route.Router,
			Link:     route.Link,
			Metric:   route.Metric,
			Metadata: route.Metadata,
		}
		respRoutes = append(respRoutes, respRoute)
	}

	resp.Routes = respRoutes

	return nil
}
