package events

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	filemanagementv1 "github.com/zchelalo/neuraclinic-file-management/gen/go/file_management/v1"
	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/ports"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NoopPublisher struct{}

func NewNoopPublisher() *NoopPublisher {
	return &NoopPublisher{}
}

func (p *NoopPublisher) PublishFileStatusChanged(_ context.Context, _ ports.FileStatusChangedEvent) error {
	return nil
}

func (p *NoopPublisher) Close() error {
	return nil
}

type RabbitPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewRabbitPublisher(url, exchange string) (*RabbitPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare rabbitmq exchange: %w", err)
	}

	return &RabbitPublisher{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
	}, nil
}

func (p *RabbitPublisher) PublishFileStatusChanged(ctx context.Context, event ports.FileStatusChangedEvent) error {
	body, err := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}.Marshal(&filemanagementv1.FileStatusChangedEvent{
		EventId:       event.EventID.String(),
		FileId:        event.FileID.String(),
		ServiceOrigin: event.ServiceOrigin,
		Status:        event.Status,
		OccurredAt:    timestamppb.New(event.OccurredAt),
		RequestId:     event.RequestID,
		TraceId:       event.TraceID,
	})
	if err != nil {
		return fmt.Errorf("marshal file status changed event: %w", err)
	}

	return p.channel.PublishWithContext(ctx, p.exchange, routingKey(event.ServiceOrigin), false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func (p *RabbitPublisher) Close() error {
	var err error
	if p.channel != nil {
		err = p.channel.Close()
	}
	if p.conn != nil {
		if closeErr := p.conn.Close(); err == nil {
			err = closeErr
		}
	}
	return err
}

func routingKey(serviceOrigin string) string {
	return "file." + serviceOrigin + ".status_changed.v1"
}

var _ ports.EventPublisher = (*NoopPublisher)(nil)
var _ ports.EventPublisher = (*RabbitPublisher)(nil)
