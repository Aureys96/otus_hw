package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/app"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/server/config"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
)

type Server struct {
	UnimplementedEventsServiceServer
	app    app.App
	cfg    config.ServerConfig
	srv    *grpc.Server
	logger *zap.Logger
}

func NewServer(app app.App, cfg config.ServerConfig, logger *zap.Logger) *Server {
	return &Server{
		app:    app,
		cfg:    cfg,
		srv:    grpc.NewServer(grpc.UnaryInterceptor(logRequest(logger))),
		logger: logger,
	}
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
	if err != nil {
		s.logger.Error("start grpc server", zap.Error(err))
		return err
	}
	RegisterEventsServiceServer(s.srv, s)
	s.logger.Info("starting grpc server",
		zap.String("host", s.cfg.Host),
		zap.Int("port", s.cfg.Port))

	return s.srv.Serve(lsn)
}

func (s *Server) Stop() error {
	s.srv.GracefulStop()
	return nil
}

func (s *Server) CreateEvent(ctx context.Context, request *CreateEventRequest) (*EventResponse, error) {
	event, err := createEventFromRequest(request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.app.CreateEvent(ctx, *event)
	if err != nil {
		if errors.Is(err, storage.ErrEventBooked) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return createEventResponse(*event)
}

func (s *Server) GetEvent(ctx context.Context, request *GetEventRequest) (*EventResponse, error) {
	event, err := s.app.Get(ctx, int(request.ID))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return createEventResponse(event)
}

func (s *Server) UpdateEvent(ctx context.Context, request *UpdateEventRequest) (*EventResponse, error) {
	event, err := updateEventFromRequest(request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.app.Update(ctx, int(request.ID), *event)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return createEventResponse(*event)
}

func (s *Server) DeleteEvent(ctx context.Context, request *DeleteEventRequest) (*emptypb.Empty, error) {
	s.app.Delete(ctx, int(request.ID))

	return &empty.Empty{}, nil
}

func (s *Server) ListEvents(ctx context.Context, request *EventListRequest) (*EventList, error) {
	var t timestamppb.Timestamp = *request.StartTime
	_, err := s.app.ListEvents(ctx, timestamppb.Timestamp.AsTime( t), timestamppb.Timestamp.AsTime(*request.EndTime))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Server) mustEmbedUnimplementedEventsServiceServer() {
	//TODO implement me
	panic("implement me")
}

func createEventFromRequest(from *CreateEventRequest) (*storage.Event, error) {
	startTime, err := ptypes.Timestamp(from.StartTime)
	if err != nil {
		return nil, err
	}
	duration, err := ptypes.Duration(from.Duration)
	if err != nil {
		return nil, err
	}
	return &storage.Event{
		Title:     from.Title,
		StartTime: startTime,
		Duration:  duration,
	}, nil
}

func updateEventFromRequest(from *UpdateEventRequest) (*storage.Event, error) {
	startTime, err := ptypes.Timestamp(from.StartTime)
	if err != nil {
		return nil, err
	}
	duration, err := ptypes.Duration(from.Duration)
	if err != nil {
		return nil, err
	}
	return &storage.Event{
		ID:        int(from.ID),
		Title:     from.Title,
		StartTime: startTime,
		Duration:  duration,
	}, nil
}

func createEventResponse(event storage.Event) (*EventResponse, error) {
	startTime, err := ptypes.TimestampProto(event.StartTime)
	if err != nil {
		return nil, err
	}

	return &EventResponse{
		ID:        int64(event.ID),
		Title:     event.Title,
		StartTime: startTime,
		Duration:  ptypes.DurationProto(event.Duration),
	}, nil
}
