package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed a_26.sql
	addExternalUserIDToLogout string
)

type UserSessionsAddExternalUserIDToLogout struct {
	dbClient *database.DB
}

func (mig *UserSessionsAddExternalUserIDToLogout) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addExternalUserIDToLogout)
	return err
}

func (mig *UserSessionsAddExternalUserIDToLogout) String() string {
	return "26_user_sessions_add_external_user_id_to_logout"
}
