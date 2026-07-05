package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	orderv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/openapi/order/v1"
	paymentv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/payment/v1"
)

var openAPIToGRPCMethod = map[orderv1.PayOrderRequestPaymentMethod]paymentv1.PaymentMethod{
	orderv1.PayOrderRequestPaymentMethodCARD:          paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	orderv1.PayOrderRequestPaymentMethodSBP:           paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
	orderv1.PayOrderRequestPaymentMethodCREDITCARD:    paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	orderv1.PayOrderRequestPaymentMethodINVESTORMONEY: paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
}

func StrToMethod(str string) (paymentv1.PaymentMethod, error) {
	strMethod, ok := openAPIToGRPCMethod[orderv1.PayOrderRequestPaymentMethod(str)]
	if !ok {
		return -1, fmt.Errorf("payment method is notFound")
	}
	return strMethod, nil
}

// ToGetOrderResponse конвертирует доменную модель Order в API-ответ GetOrderResponse.
func ToGetOrderResponse(order *model.Order) *orderv1.GetOrderResponse {
	if order == nil {
		return nil
	}

	resp := &orderv1.GetOrderResponse{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     convertStatusToAPI(order.Status),
	}

	// Опциональное поле: TransactionUUID
	if order.TransactionUUID != nil {
		resp.TransactionUUID = orderv1.NewOptNilString(*order.TransactionUUID)
	}

	// Опциональное поле: PaymentMethod
	if order.PaymentMethod != nil {
		method := orderv1.GetOrderResponsePaymentMethod(*order.PaymentMethod)
		resp.PaymentMethod = orderv1.NewOptNilGetOrderResponsePaymentMethod(method)
	}

	return resp
}

// ToCreateOrderResponse конвертирует доменную модель Order в API-ответ CreateOrderResponse.
func ToCreateOrderResponse(order *model.Order) *orderv1.CreateOrderResponse {
	if order == nil {
		return nil
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderv1.NewOptUUID(uuid.MustParse(order.OrderUUID)),
		TotalPrice: orderv1.NewOptFloat64(order.TotalPrice),
	}
}

// ToPayOrderResponse конвертирует transaction UUID в API-ответ PayOrderResponse.
func ToPayOrderResponse(PtrTransactionUUID *string) *orderv1.PayOrderResponse {
	transactionUUID := *PtrTransactionUUID
	return &orderv1.PayOrderResponse{
		TransactionUUID: orderv1.NewOptString(transactionUUID),
	}
}

// convertStatusToAPI конвертирует доменный статус в API-статус.
func convertStatusToAPI(status model.OrderStatus) orderv1.GetOrderResponseStatus {
	switch status {
	case model.StatusPendingPayment:
		return orderv1.GetOrderResponseStatusPENDINGPAYMENT
	case model.StatusPaid:
		return orderv1.GetOrderResponseStatusPAID
	case model.StatusCancelled:
		return orderv1.GetOrderResponseStatusCANCELLED
	case model.StatusAssembled:
		return orderv1.GetOrderResponseStatusASSEMBLED

	default:
		return orderv1.GetOrderResponseStatusPENDINGPAYMENT
	}
}
