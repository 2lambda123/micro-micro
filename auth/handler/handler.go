package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/auth/token"
	"github.com/micro/go-micro/v2/auth/token/basic"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	memStore "github.com/micro/go-micro/v2/store/memory"
)

var joinKey = ":"

// Handler processes RPC calls
type Handler struct {
	Options        auth.Options
	SecretProvider token.Provider
	TokenProvider  token.Provider
}

// Init the auth
func (h *Handler) Init(opts ...auth.Option) {
	for _, o := range opts {
		o(&h.Options)
	}

	// use the default store as a fallback
	if h.Options.Store == nil {
		h.Options.Store = store.DefaultStore
	}

	// noop will not work for auth
	if h.Options.Store.String() == "noop" {
		h.Options.Store = memStore.NewStore()
	}

	if h.TokenProvider == nil {
		h.TokenProvider = basic.NewTokenProvider(token.WithStore(h.Options.Store))
	}
	if h.SecretProvider == nil {
		h.SecretProvider = basic.NewTokenProvider(token.WithStore(h.Options.Store))
	}
}

// Generate an account
func (h *Handler) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	// Generate a long-lived secret
	secretOpts := []token.GenerateOption{
		token.WithExpiry(time.Duration(req.SecretExpiry) * time.Second),
		token.WithMetadata(req.Metadata),
		token.WithRoles(req.Roles...),
	}
	secret, err := h.SecretProvider.Generate(req.Id, secretOpts...)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to generate secret: %v", err)
	}

	// set the account
	rsp.Account = &pb.Account{
		Id:       req.Id,
		Metadata: req.Metadata,
		Roles:    req.Roles,
		Secret:   serializeToken(secret),
	}

	return nil
}

// Grant a role access to a resource
func (h *Handler) Grant(ctx context.Context, req *pb.GrantRequest, rsp *pb.GrantResponse) error {
	if req.Resource == nil {
		return errors.BadRequest("go.micro.auth", "Resource missing")
	}

	// Construct the key
	comps := []string{req.Resource.Type, req.Resource.Name, req.Resource.Endpoint, req.Role}
	key := strings.Join(comps, joinKey)

	// Encode the rule
	bytes, err := json.Marshal(pb.Rule{Role: req.Role, Resource: req.Resource})
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to marshal rule: %v", err)
	}

	// Write to the store
	if err := h.Options.Store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to write to the store: %v", err)
	}

	return nil
}

// Revoke a roles access to a resource
func (h *Handler) Revoke(ctx context.Context, req *pb.RevokeRequest, rsp *pb.RevokeResponse) error {
	if req.Resource == nil {
		return errors.BadRequest("go.micro.auth", "Resource missing")
	}

	// Construct the key
	comps := []string{req.Resource.Type, req.Resource.Name, req.Resource.Endpoint, req.Role}
	key := strings.Join(comps, joinKey)

	// Delete the rule
	err := h.Options.Store.Delete(key)
	if err == store.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Rule not found")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to delete key from store: %v", err)
	}

	return nil
}

// Inspect a token and retrieve the account
func (h *Handler) Inspect(ctx context.Context, req *pb.InspectRequest, rsp *pb.InspectResponse) error {
	tok, err := h.TokenProvider.Inspect(req.Token)
	if err == token.ErrInvalidToken || err == token.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Invalid token")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to inspect token: %v", err)
	}

	rsp.Account = &pb.Account{
		Id:       tok.Subject,
		Roles:    tok.Roles,
		Metadata: tok.Metadata,
	}
	return nil
}

// Refresh a token using a secret
func (h *Handler) Refresh(ctx context.Context, req *pb.RefreshRequest, rsp *pb.RefreshResponse) error {
	sec, err := h.SecretProvider.Inspect(req.Secret)
	if err == token.ErrInvalidToken || err == token.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Invalid token")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to inspect secret: %v", err)
	}

	tok, err := h.TokenProvider.Generate(sec.Subject,
		token.WithExpiry(time.Duration(req.TokenExpiry)*time.Second),
		token.WithMetadata(sec.Metadata),
		token.WithRoles(sec.Roles...),
	)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to generate token: %v", err)
	}

	rsp.Token = serializeToken(tok)
	return nil
}

// ListRules returns all the rules
func (h *Handler) ListRules(ctx context.Context, req *pb.ListRulesRequest, rsp *pb.ListRulesResponse) error {
	// get the records from the store
	recs, err := h.Options.Store.Read("", store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to read from store: %v", err)
	}

	// unmarshal the records
	rsp.Rules = make([]*pb.Rule, 0, len(recs))
	for _, rec := range recs {
		var r *pb.Rule
		if err := json.Unmarshal(rec.Value, &r); err != nil {
			return errors.InternalServerError("go.micro.auth", "Error to unmarshaling json: %v. Value: %v", err, string(rec.Value))
		}
		r.Id = rec.Key
		rsp.Rules = append(rsp.Rules, r)
	}

	return nil
}

func serializeAccount(a *auth.Account) *pb.Account {
	var secret *pb.Token
	if a.Secret != nil {
		secret = serializeToken(a.Secret)
	}

	return &pb.Account{
		Id:       a.ID,
		Roles:    a.Roles,
		Metadata: a.Metadata,
		Secret:   secret,
	}
}

func serializeToken(t *auth.Token) *pb.Token {
	return &pb.Token{
		Token:    t.Token,
		Type:     t.Type,
		Created:  t.Created.Unix(),
		Expiry:   t.Expiry.Unix(),
		Subject:  t.Subject,
		Roles:    t.Roles,
		Metadata: t.Metadata,
	}
}
