package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OrgOIDCSettingsWriteModel struct {
	eventstore.WriteModel

	AccessTokenLifetime        time.Duration
	IdTokenLifetime            time.Duration
	RefreshTokenIdleExpiration time.Duration
	RefreshTokenExpiration     time.Duration
	State                      domain.OIDCSettingsState
}

func NewOrgOIDCSettingsWriteModel(orgID string) *OrgOIDCSettingsWriteModel {
	return &OrgOIDCSettingsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   orgID,
			ResourceOwner: orgID,
		},
	}
}

func (wm *OrgOIDCSettingsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.OIDCSettingsAddedEvent:
			wm.AccessTokenLifetime = e.AccessTokenLifetime
			wm.IdTokenLifetime = e.IdTokenLifetime
			wm.RefreshTokenIdleExpiration = e.RefreshTokenIdleExpiration
			wm.RefreshTokenExpiration = e.RefreshTokenExpiration
			wm.State = domain.OIDCSettingsStateActive
		case *org.OIDCSettingsChangedEvent:
			if e.AccessTokenLifetime != nil {
				wm.AccessTokenLifetime = *e.AccessTokenLifetime
			}
			if e.IdTokenLifetime != nil {
				wm.IdTokenLifetime = *e.IdTokenLifetime
			}
			if e.RefreshTokenIdleExpiration != nil {
				wm.RefreshTokenIdleExpiration = *e.RefreshTokenIdleExpiration
			}
			if e.RefreshTokenExpiration != nil {
				wm.RefreshTokenExpiration = *e.RefreshTokenExpiration
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OrgOIDCSettingsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.OIDCSettingsAddedEventType,
			org.OIDCSettingsChangedEventType).
		Builder()
}

func (wm *OrgOIDCSettingsWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	accessTokenLifetime,
	idTokenLifetime,
	refreshTokenIdleExpiration,
	refreshTokenExpiration time.Duration,
) (*org.OIDCSettingsChangedEvent, bool, error) {
	changes := make([]org.OIDCSettingsChanges, 0, 4)
	var err error

	if wm.AccessTokenLifetime != accessTokenLifetime {
		changes = append(changes, org.ChangeOIDCSettingsAccessTokenLifetime(accessTokenLifetime))
	}
	if wm.IdTokenLifetime != idTokenLifetime {
		changes = append(changes, org.ChangeOIDCSettingsIdTokenLifetime(idTokenLifetime))
	}
	if wm.RefreshTokenIdleExpiration != refreshTokenIdleExpiration {
		changes = append(changes, org.ChangeOIDCSettingsRefreshTokenIdleExpiration(refreshTokenIdleExpiration))
	}
	if wm.RefreshTokenExpiration != refreshTokenExpiration {
		changes = append(changes, org.ChangeOIDCSettingsRefreshTokenExpiration(refreshTokenExpiration))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := org.NewOIDCSettingsChangeEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
