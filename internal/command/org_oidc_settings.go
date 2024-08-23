package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) prepareAddOrgOIDCSettings(a *org.Aggregate, accessTokenLifetime, idTokenLifetime, refreshTokenIdleExpiration, refreshTokenExpiration time.Duration) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if accessTokenLifetime == time.Duration(0) ||
			idTokenLifetime == time.Duration(0) ||
			refreshTokenIdleExpiration == time.Duration(0) ||
			refreshTokenExpiration == time.Duration(0) {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-10s82j", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := c.getOrgOIDCSettingsWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			if writeModel.State == domain.OIDCSettingsStateActive {
				return nil, zerrors.ThrowAlreadyExists(nil, "INST-0aaj1o", "Errors.OIDCSettings.AlreadyExists")
			}
			return []eventstore.Command{
				org.NewOIDCSettingsAddedEvent(
					ctx,
					&a.Aggregate,
					accessTokenLifetime,
					idTokenLifetime,
					refreshTokenIdleExpiration,
					refreshTokenExpiration,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgOIDCSettings(a *org.Aggregate, accessTokenLifetime, idTokenLifetime, refreshTokenIdleExpiration, refreshTokenExpiration time.Duration) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if accessTokenLifetime == time.Duration(0) ||
			idTokenLifetime == time.Duration(0) ||
			refreshTokenIdleExpiration == time.Duration(0) ||
			refreshTokenExpiration == time.Duration(0) {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-10sxks", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := c.getOrgOIDCSettingsWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			if writeModel.State != domain.OIDCSettingsStateActive {
				return []eventstore.Command{
					org.NewOIDCSettingsAddedEvent(
						ctx,
						&a.Aggregate,
						accessTokenLifetime,
						idTokenLifetime,
						refreshTokenIdleExpiration,
						refreshTokenExpiration,
					),
				}, nil
			}
			changedEvent, hasChanged, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				accessTokenLifetime,
				idTokenLifetime,
				refreshTokenIdleExpiration,
				refreshTokenExpiration,
			)
			if err != nil {
				return nil, err
			}
			if !hasChanged {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-0pk2nu", "Errors.NoChangesFound")
			}
			return []eventstore.Command{
				changedEvent,
			}, nil
		}, nil
	}
}

func (c *Commands) AddOrgOIDCSettings(ctx context.Context, settings *domain.OIDCSettings) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(authz.GetCtxData(ctx).OrgID)
	validation := c.prepareAddOrgOIDCSettings(orgAgg, settings.AccessTokenLifetime, settings.IdTokenLifetime, settings.RefreshTokenIdleExpiration, settings.RefreshTokenExpiration)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func (c *Commands) ChangeOrgOIDCSettings(ctx context.Context, settings *domain.OIDCSettings) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(authz.GetCtxData(ctx).OrgID)
	validation := c.prepareUpdateOrgOIDCSettings(orgAgg, settings.AccessTokenLifetime, settings.IdTokenLifetime, settings.RefreshTokenIdleExpiration, settings.RefreshTokenExpiration)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}, nil
}

func (c *Commands) getOrgOIDCSettingsWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer) (_ *OrgOIDCSettingsWriteModel, err error) {
	orgID := authz.GetCtxData(ctx).OrgID
	writeModel := NewOrgOIDCSettingsWriteModel(orgID)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
