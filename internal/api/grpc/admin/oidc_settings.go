package admin

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	grpc_util "github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetOIDCSettings(ctx context.Context, _ *admin_pb.GetOIDCSettingsRequest) (*admin_pb.GetOIDCSettingsResponse, error) {
	var err error
	var result *query.OIDCSettings
	if orgID := grpc_util.GetHeader(ctx, http.ZitadelOrgID); orgID != "" {
		result, err = s.query.OIDCSettingsByAggID(ctx, orgID)
		if strings.Contains(err.Error(), "QUERY-s9nlw") {
			// ignore NotFound error
			err = nil
		}
	}
	if result == nil {
		result, err = s.query.OIDCSettingsByAggID(ctx, authz.GetInstance(ctx).InstanceID())
	}
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetOIDCSettingsResponse{
		Settings: OIDCSettingsToPb(result),
	}, nil
}

func (s *Server) AddOIDCSettings(ctx context.Context, req *admin_pb.AddOIDCSettingsRequest) (*admin_pb.AddOIDCSettingsResponse, error) {
	var err error
	var result *domain.ObjectDetails
	if orgID := grpc_util.GetHeader(ctx, http.ZitadelOrgID); orgID != "" {
		result, err = s.command.AddOrgOIDCSettings(ctx, AddOIDCConfigToConfig(req))
	} else {
		result, err = s.command.AddOIDCSettings(ctx, AddOIDCConfigToConfig(req))
	}
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddOIDCSettingsResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) UpdateOIDCSettings(ctx context.Context, req *admin_pb.UpdateOIDCSettingsRequest) (*admin_pb.UpdateOIDCSettingsResponse, error) {
	var err error
	var result *domain.ObjectDetails
	if orgID := grpc_util.GetHeader(ctx, http.ZitadelOrgID); orgID != "" {
		result, err = s.command.ChangeOrgOIDCSettings(ctx, UpdateOIDCConfigToConfig(req))
	} else {
		result, err = s.command.ChangeOIDCSettings(ctx, UpdateOIDCConfigToConfig(req))
	}
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateOIDCSettingsResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}
