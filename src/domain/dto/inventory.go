package dto

import "crazyfarmbackend/src/constant"

type InventoryItem struct {
	Plant    constant.Plant `json:"Plant"`
	Quantity int            `json:"Quantity"`
}

type GetAllItemsResponse struct {
	Items []InventoryItem `json:"items"`
}
