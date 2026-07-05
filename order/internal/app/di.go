package app

import (
	api "github.com/Ilia9531/microservices-warehouse/order/internal/api/order/v1"
	cl "github.com/Ilia9531/microservices-warehouse/order/internal/client/grpc"
	invCl "github.com/Ilia9531/microservices-warehouse/order/internal/client/grpc/inventory/v1"
	payCl "github.com/Ilia9531/microservices-warehouse/order/internal/client/grpc/payment/v1"
	kafkaConverter "github.com/Ilia9531/microservices-warehouse/order/internal/converter/kafka"
	"github.com/Ilia9531/microservices-warehouse/order/internal/converter/kafka/decoder"
	"github.com/Ilia9531/microservices-warehouse/order/internal/migrator"
	repo "github.com/Ilia9531/microservices-warehouse/order/internal/repository/order"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/closer"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
	"context"
	"fmt"

	orderConsumer "github.com/Ilia9531/microservices-warehouse/order/internal/service/consumer/order_consumer"
	kafkaMiddleware "github.com/Ilia9531/microservices-warehouse/platform/pkg/middleware/kafka"

	"github.com/Ilia9531/microservices-warehouse/order/internal/config"
	"github.com/Ilia9531/microservices-warehouse/order/internal/repository"
	"github.com/Ilia9531/microservices-warehouse/order/internal/service"
	orderScv "github.com/Ilia9531/microservices-warehouse/order/internal/service/order"
	orderProducer "github.com/Ilia9531/microservices-warehouse/order/internal/service/producer/order_producer"
	wrappedKafka "github.com/Ilia9531/microservices-warehouse/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Ilia9531/microservices-warehouse/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/Ilia9531/microservices-warehouse/platform/pkg/kafka/producer"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	pgConn *pgx.Conn // Подключение к PostgreSQL (pgx)

	inventoryClient cl.InventoryClient // gRPC-клиент к InventoryService
	paymentClient   cl.PaymentClient   // gRPC-клиент к PaymentService

	apiHandler   *api.Api                   // HTTP API handler (OpenAPI + chi)
	orderService service.OrderService       // Бизнес-логика
	orderRepo    repository.OrderRepository // Репозиторий заказов

	orderProducerService service.OrderProducerService
	orderConsumerService service.ConsumerService

	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderAssembledConsumer      wrappedKafka.Consumer
	shipAssembledDecoder        kafkaConverter.ShipAssembledDecoder

	syncProducer      sarama.SyncProducer
	orderPaidProducer wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PostgresConn(ctx context.Context) *pgx.Conn {
	if d.pgConn == nil {
		conn, err := pgx.Connect(ctx, config.AppConfig().Postgres.DSN())
		if err != nil {
			panic(fmt.Sprintf("❌ Failed to connect to database: %v", err))
		}

		if err = conn.Ping(ctx); err != nil {
			panic(fmt.Sprintf("❌ Database ping failed: %v", err))
		}
		closer.AddNamed("Postgres client", func(ctx context.Context) error {
			return conn.Close(ctx)
		})
		logger.Info(ctx, "✅ Connected to PostgreSQL")
		d.pgConn = conn
	}
	return d.pgConn
}

func (d *diContainer) RunMigrations(ctx context.Context) {
	cfg := config.AppConfig().Postgres

	dbForMigrations := stdlib.OpenDB(*d.PostgresConn(ctx).Config().Copy())
	defer dbForMigrations.Close()

	mig := migrator.NewMigrator(dbForMigrations, cfg.MigrationsDir())
	if err := mig.Up(); err != nil {
		panic(fmt.Sprintf("failed to run migrations: %v", err))
	}
	logger.Info(ctx, "✅ Migrations completed")
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {

	if d.orderRepo == nil {
		d.orderRepo = repo.NewRepository(d.PostgresConn(ctx))
	}
	logger.Info(ctx, "✅ Repository инициализирован")

	return d.orderRepo
}
func (d *diContainer) InventoryClient(ctx context.Context) cl.InventoryClient {
	if d.inventoryClient == nil {
		cfg := config.AppConfig().InvGRPC

		conn, err := grpc.NewClient(
			cfg.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create inventory grpc connection: %v", err))
		}

		client := invCl.NewClient(conn)

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryClient = client
		logger.Info(ctx, "✅ InvgRPC-клиент инициализирован")

	}
	return d.inventoryClient
}
func (d *diContainer) PaymentClient(ctx context.Context) cl.PaymentClient {
	if d.paymentClient == nil {
		cfg := config.AppConfig().PayGRPC

		conn, err := grpc.NewClient(
			cfg.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create payment grpc connection: %v", err))
		}

		client := payCl.NewPaymentClient(conn)

		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentClient = client
		logger.Info(ctx, "✅ PayRPC-клиент инициализирован")
	}
	return d.paymentClient
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderPaidProducer())
	}
	return d.orderProducerService
}
func (d *diContainer) OrderConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(
			d.OrderAssembledConsumer(),
			d.ShipAssembledDecoder(),
			d.OrderRepository(ctx),
		)
	}
	return d.orderConsumerService
}
func (d *diContainer) OrderPaidProducer() wrappedKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
		logger.Info(context.Background(), "✅ OrderPaid Kafka producer initialized")
	}
	return d.orderPaidProducer
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		cfg := config.AppConfig().OrderPaidProducer.Config()
		p, err := sarama.NewSyncProducer(config.AppConfig().Kafka.Brokers(), cfg)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer (OrderPaid)", func(ctx context.Context) error {
			return p.Close()
		})
		d.syncProducer = p
	}
	return d.syncProducer
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		cfg := config.AppConfig().OrderAssembledConsumer.Config()
		cg, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			cfg,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group (OrderAssembled)", func(ctx context.Context) error {
			return cg.Close()
		})
		d.orderAssembledConsumerGroup = cg
	}
	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{config.AppConfig().OrderAssembledConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
		logger.Info(context.Background(), "✅ OrderAssembled Kafka consumer initialized")
	}
	return d.orderAssembledConsumer
}

func (d *diContainer) ShipAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewShipAssembledDecoder()
	}
	return d.shipAssembledDecoder
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderScv.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
			d.OrderProducerService(),
		)
	}
	return d.orderService
}

func (d *diContainer) APIHandler(ctx context.Context) *api.Api {
	if d.apiHandler == nil {
		d.apiHandler = api.NewApi(d.OrderService(ctx))
	}
	return d.apiHandler
}
