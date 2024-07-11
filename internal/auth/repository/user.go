package repository

import (
	"context"

	"github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type UserRepository interface {
	UserSessionsByAgentID(ctx context.Context, agentID string) ([]*model.UserSessionView, error)
	UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error)
}
