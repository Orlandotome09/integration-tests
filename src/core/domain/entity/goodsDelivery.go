package entity

type GoodsDelivery struct {
	AverageDeliveryDays   int            `json:"average_delivery_days"`
	ShippingMethods       ShippingMethod `json:"shipping_methods"`
	Insurance             bool           `json:"insurance"`
	TrackingCodeAvailable bool           `json:"tracking_code_available"`
}
