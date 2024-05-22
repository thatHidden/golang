package models

type CarBaseInputDTO struct {
	Id           uint64  `json:"id"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model"`
	Generation   string  `json:"generation"`
	BodyStyle    string  `json:"body_style"`
	Transmission string  `json:"transmission"`
	Drivetrain   string  `json:"drivetrain"`
	EngineSize   float64 `json:"engine_size"`
	EngineType   string  `json:"engine_type"`
	HorsePower   uint16  `json:"horse_power"`
	Torque       uint16  `json:"torque"`
	BuildFrom    uint16  `json:"build_from"`
	BuildTo      uint16  `json:"build_to"`
}

type CarBaseOutputDTO struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Generation string `json:"generation"`
	BodyStyle  string `json:"body_style"`
	PowerTrain struct {
		Transmission string `json:"transmission"`
		Drivetrain   string `json:"drivetrain"`
		Engine       struct {
			EngineSize string `json:"size"`
			EngineType string `json:"type"`
			HorsePower string `json:"power"`
			Torque     string `json:"torque"`
		} `json:"engine"`
	} `json:"power_train"`
	BuildYears string `json:"build_years"`
}

//func newCarBaseDTO(u *domain.BaseCar) {
//	//returns DTO based on domain
//}
//
//func newCarBaseDomain(dto *CarBaseInputDTO) {
//	//returns domain based on dto
//}

//// не нужно?
//type CarBaseThinOutputDTO struct {
//	Name       string `json:"name"`
//	Generation string `json:"generation"`
//	BodyStyle  string `json:"body_style"`
//}
