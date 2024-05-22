package models

type CarInputDTO struct {
	Id            uint64 `json:"id"`
	BaseCarID     uint64 `json:"base_car_id"`
	CarModsID     uint64 `json:"car_mods_id"`
	Mileage       uint32 `json:"mileage"`
	Vin           string `json:"vin"`
	Location      string `json:"location"`
	ExteriorColor string `json:"exterior_color"`
	InteriorColor string `json:"interior_color"`
}

type CarOutputDTO struct {
	CarMods  CarModsOutputDTO `json:"mods"`
	CarBase  CarBaseOutputDTO `json:"base"`
	Mileage  uint32           `json:"mileage"`
	Vin      string           `json:"vin"`
	Location string           `json:"location"`
	Color    struct {
		ExteriorColor string `json:"exterior"`
		InteriorColor string `json:"interior"`
	} `json:"color"`
}
