package entity

// LocationUpdate represents a request to update inventory storage location
type LocationUpdate struct {
	BinLocation string `json:"bin_location"`
	ShelfNumber string `json:"shelf_number"`
	ZoneCode    string `json:"zone_code"`
}
