package server

import (
	"context"

	"github.com/micro/go-micro/v3/network"
	"github.com/micro/go-micro/v3/network/mucp"
	"github.com/micro/go-micro/v3/router"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/micro/v3/service/network/proto"
	"github.com/micro/micro/v3/service/network/util"
	pbRtr "github.com/micro/micro/v3/service/router/proto"
)

// Network implements network handler
type Network struct {
	Network network.Network
}

func flatten(n network.Node, visited map[string]bool) []network.Node {
	// if node is nil runaway
	if n == nil {
		return nil
	}

	// set visisted
	if visited == nil {
		visited = make(map[string]bool)
	}

	// create new list of nodes
	//nolint:prealloc
	var nodes []network.Node

	// check if already visited
	if !visited[n.Id()] {
		// append the current node
		nodes = append(nodes, n)
	}

	// set to visited
	visited[n.Id()] = true

	// visit the list of peers
	for _, node := range n.Peers() {
		nodes = append(nodes, flatten(node, visited)...)
	}

	return nodes
}

func (n *Network) Connect(ctx context.Context, req *pb.ConnectRequest, resp *pb.ConnectResponse) error {
	if len(req.Nodes) == 0 {
		return nil
	}

	// get list of existing nodes
	nodes := n.Network.Options().Nodes

	// generate a node map
	nodeMap := make(map[string]bool)

	for _, node := range nodes {
		nodeMap[node] = true
	}

	for _, node := range req.Nodes {
		// TODO: we may have been provided a network only
		// so process anad resolve node.Network
		if len(node.Address) == 0 {
			continue
		}

		// already exists
		if _, ok := nodeMap[node.Address]; ok {
			continue
		}

		nodeMap[node.Address] = true
		nodes = append(nodes, node.Address)
	}

	log.Infof("Network.Connect setting peers: %v", nodes)

	// reinitialise the peers
	n.Network.Init(
		network.Nodes(nodes...),
	)

	// call the connect method
	n.Network.Connect()

	return nil
}

// Nodes returns the list of nodes
func (n *Network) Nodes(ctx context.Context, req *pb.NodesRequest, resp *pb.NodesResponse) error {
	// root node
	nodes := map[string]network.Node{}

	// get peers encoded into protobuf
	peers := flatten(n.Network, nil)

	// walk all the peers
	for _, peer := range peers {
		if peer == nil {
			continue
		}
		if _, ok := nodes[peer.Id()]; ok {
			continue
		}

		// add to visited list
		nodes[n.Network.Id()] = peer

		resp.Nodes = append(resp.Nodes, &pb.Node{
			Id:      peer.Id(),
			Address: peer.Address(),
		})
	}

	return nil
}

// Graph returns the network graph from this root node
func (n *Network) Graph(ctx context.Context, req *pb.GraphRequest, resp *pb.GraphResponse) error {
	depth := uint(req.Depth)
	if depth <= 0 || depth > mucp.MaxDepth {
		depth = mucp.MaxDepth
	}

	// get peers encoded into protobuf
	peers := util.PeersToProto(n.Network, depth)

	// set the root node
	resp.Root = peers

	return nil
}

// Routes returns a list of routing table routes
func (n *Network) Routes(ctx context.Context, req *pb.RoutesRequest, resp *pb.RoutesResponse) error {
	// build query

	var qOpts []router.QueryOption

	if q := req.Query; q != nil {
		if len(q.Service) > 0 {
			qOpts = append(qOpts, router.QueryService(q.Service))
		}
		if len(q.Address) > 0 {
			qOpts = append(qOpts, router.QueryAddress(q.Address))
		}
		if len(q.Gateway) > 0 {
			qOpts = append(qOpts, router.QueryGateway(q.Gateway))
		}
		if len(q.Router) > 0 {
			qOpts = append(qOpts, router.QueryRouter(q.Router))
		}
		if len(q.Network) > 0 {
			qOpts = append(qOpts, router.QueryNetwork(q.Network))
		}
	}

	routes, err := n.Network.Options().Router.Table().Query(qOpts...)
	if err != nil {
		return errors.InternalServerError("go.micro.network", "failed to list routes: %s", err)
	}

	respRoutes := make([]*pbRtr.Route, 0, len(routes))
	for _, route := range routes {
		respRoute := &pbRtr.Route{
			Service: route.Service,
			Address: route.Address,
			Gateway: route.Gateway,
			Network: route.Network,
			Router:  route.Router,
			Link:    route.Link,
			Metric:  int64(route.Metric),
		}
		respRoutes = append(respRoutes, respRoute)
	}

	resp.Routes = respRoutes

	return nil
}

// Services returns a list of services based on the routing table
func (n *Network) Services(ctx context.Context, req *pb.ServicesRequest, resp *pb.ServicesResponse) error {
	routes, err := n.Network.Options().Router.Table().List()
	if err != nil {
		return errors.InternalServerError("go.micro.network", "failed to list services: %s", err)
	}

	services := make(map[string]bool)

	for _, route := range routes {
		if route.Service == "*" {
			continue
		}

		if _, ok := services[route.Service]; ok {
			continue
		}
		services[route.Service] = true
		resp.Services = append(resp.Services, route.Service)
	}

	return nil
}
