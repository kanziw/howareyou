package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/kanziw/go-slack"
	"github.com/kanziw/howareyou/model"
)

type Service interface {
	StartHowAreYou(ctx context.Context, channel, userGroup string) error
}

type DefaultService struct {
	db  *sql.DB
	api *slack.Client
}

func (s *DefaultService) StartHowAreYou(ctx context.Context, channel, userGroup string) error {
	// TODO: Upsert Schedule into DB
	ug, err := model.UserGroups(qm.Where(model.UserGroupColumns.UserGroupSlackID, userGroup)).One(ctx, s.db)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return errors.Wrap(err, "find user group")
	}
	if ug == nil {
		ug = &model.UserGroup{
			UserGroupSlackID: userGroup,
			ChannelSlackID:   "",
		}
		if err := ug.Insert(ctx, s.db, boil.Infer()); err != nil {
			return errors.Wrap(err, "insert user group")
		}
	}

	if err := slack.SendMessage(ctx, s.api, channel, "Let's start HowAreYou with "+userGroup); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *DefaultService) SendHowAreYou(ctx context.Context, channel, userGroup string) error {
	if err := slack.SendMessage(
		ctx,
		s.api,
		channel,
		fmt.Sprintf("Hi %s :) How are you today?", userGroup),
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func New(db *sql.DB, api *slack.Client) Service {
	return &DefaultService{
		db:  db,
		api: api,
	}
}
