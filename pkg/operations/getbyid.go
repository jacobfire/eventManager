package operations

import (
	"calendar/internal/app/model"
	"calendar/internal/app/storage"
	"calendar/pkg/api"
	"context"
	"fmt"
)

type GRPCServer struct {}

func (s *GRPCServer) GetById(ctx context.Context, req *api.GetByIdRequest) (*api.GetByIdResponse, error) {
	event := model.Event{}
	storage.GetById(&event, req.Id)

	fmt.Printf("%+v", event)

	return &api.GetByIdResponse{
		Id: int32(event.ID),
		Title: event.Title,
		Description: event.Description,
		Time: event.Time,
		Timezone: event.Timezone,
	}, nil
}

func (s *GRPCServer) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	//event := model.Event{}

	event := model.ExtendedEvent{}
	event.ID = 0
	event.Title = req.Title
	event.Description = req.Description
	event.Time = req.Time
	event.Timezone = req.Timezone
	event.Duration = req.Duration
	event.Notes = req.Notes
	storage.Create(&event)

	fmt.Printf("%+v", event)

	return &api.CreateResponse{
		Id: int32(event.ID),
		Title: event.Title,
		Description: event.Description,
		Time: event.Time,
		Timezone: event.Timezone,
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return &api.LoginResponse{
		Token: "token",
		Error: "token",
	}, nil
}

func (s *GRPCServer) All(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	var modelEvents []*model.ExtendedEvent
	generatedEvents := []*api.Event{}
	storage.All(&modelEvents)
	fmt.Println(modelEvents)

	for _, v := range modelEvents {
		generatedEvents = append(generatedEvents, &api.Event{
			Id: int32(v.ID),
			Title: v.Title,
			Description: v.Description,
			Time: v.Time,
			Timezone: v.Timezone,

		})
	}

	return &api.ListResponse{
		Events: generatedEvents,
	}, nil
}

func (s *GRPCServer) Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponse, error) {
	changes := make(map[string]interface{})
	if req.Title != "" {
		changes["title"] = req.Title
	}
	if req.Description != "" {
		changes["description"] = req.Description
	}

	if req.Time != "" {
		changes["time"] = req.Time
	}

	if req.Timezone != "" {
		changes["timezone"] = req.Timezone
	}

	if req.Duration != "" {
		changes["duration"] = req.Duration
	}

	if req.Notes != "" {
		changes["notes"] = req.Notes
	}

	storage.Update(changes, int(req.Id))

	return &api.UpdateResponse{
		Status: "status",
		Error: "no error",
	}, nil
}

func (s *GRPCServer) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	event := model.ExtendedEvent{}
	if err := storage.Delete(event, req.Id); err != nil {
		return &api.DeleteResponse{
			Status: "error",
			Error: err.Error(),
		}, nil
	}
	return &api.DeleteResponse{
		Status: "ok",
		Error: "no error",
	}, nil
}
