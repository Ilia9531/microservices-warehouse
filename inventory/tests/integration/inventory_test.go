//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	invV1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		invClient invV1.InventoryServiceClient
	)

	BeforeEach(func() {

		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC-клиент к запущенному в тесте Inventory-сервису
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		invClient = invV1.NewInventoryServiceClient(conn)

	})

	AfterEach(func() {
		// Чистим коллекцию после каждого теста для полной изоляции
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()

	})

	Describe("GetPart", func() {
		var partUUID string

		BeforeEach(func() {
			// Вставляем тестовую деталь напрямую в MongoDB через хелпер
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := invClient.GetPart(ctx, &invV1.GetPartRequest{
				Uuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().GetUuid()).To(Equal(partUUID))
			Expect(resp.GetPart().GetName()).ToNot(BeEmpty())
			Expect(resp.GetPart().GetPrice()).To(BeNumerically(">", 0))
		})

		It("должен возвращать ошибку для несуществующего UUID", func() {
			_, err := invClient.GetPart(ctx, &invV1.GetPartRequest{
				Uuid: "00000000-0000-0000-0000-000000000000",
			})

			// gRPC вернёт ошибку, так как детали нет в БД
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ListParts", func() {
		It("должен возвращать пустой список при отсутствии деталей", func() {
			// Коллекция уже очищена в AfterEach
			resp, err := invClient.ListParts(ctx, &invV1.ListPartsRequest{
				Filter: &invV1.PartsFilter{},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(BeEmpty())
		})

		It("должен возвращать все детали при пустом фильтре", func() {
			// Вставляем две разные детали
			_, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred())
			_, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred())

			resp, err := invClient.ListParts(ctx, &invV1.ListPartsRequest{
				Filter: &invV1.PartsFilter{},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(HaveLen(2))
		})

		It("должен корректно фильтровать детали по имени", func() {
			// Вставляем деталь с конкретным именем
			testPart := env.GetTestPart()
			testPart.Name = "TestEngineV1"
			_, err := env.InsertTestPartWithData(ctx, testPart)
			Expect(err).ToNot(HaveOccurred())

			// Вставляем ещё одну с другим именем
			testPart2 := env.GetTestPart()
			testPart2.Name = "OtherPartV2"
			_, err = env.InsertTestPartWithData(ctx, testPart2)
			Expect(err).ToNot(HaveOccurred())

			// Фильтруем только по первому имени
			resp, err := invClient.ListParts(ctx, &invV1.ListPartsRequest{
				Filter: &invV1.PartsFilter{
					Names: []string{"TestEngineV1"},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(HaveLen(1))
			Expect(resp.GetParts()[0].GetName()).To(Equal("TestEngineV1"))
		})
	})
})
