package notification

import (
	"fmt"
	"testing"
)

func TestNotification(t *testing.T) {
	data := `{
    "id": "6698f632-5650-413d-88ac-ad04aa6fde31",
    "entityId": "a7de1d12-dae2-420f-8432-3a470e2c3cef",
    "entityType": "business",
    "bankName": "bbva",
    "sourceId": "NO-572ea412-4b42-40f2-ae55-1357fa8c43e1",
    "type": "transaction",
    "action": "posted",
    "attribute": null,
    "version": "1.0.0",
    "created": "2020-06-25T18:07:10Z",
    "data": {
        "type": "transfer",
        "amount": 2.2,
        "cardId": null,
        "bankName": "bbva",
        "codeType": "creditPosted",
        "currency": "usd",
        "accountId": "AC-aee881b2-2799-4908-b6f7-6ffa8a4940fb",
        "postedBalance": 10776.37,
        "cardTransaction": null,
        "holdTransaction": null,
        "transactionDate": "2020-06-25T18:07:08Z",
        "bankTransactionId": "202006251307082341243DL03DL0",
        "bankMoneyTransferId": "MM-377df4ed-d814-433e-a3df-80e67c7a2d06",
        "bankTransactionDesc": "0480 Transfer From: Wise User Account: ****4872 MM-377df4ed-d814-433e-a3df-80e67c7a2d06"
    }
}`
	err := HandleNotification(&data)
	fmt.Printf("\n %+v \n", err)
}