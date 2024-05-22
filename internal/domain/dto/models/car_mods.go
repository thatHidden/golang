package models

type CarModsInputDTO struct {
	Id           uint64  `json:"id"` //удоли
	BodyStyle    string  `json:"body_style"`
	Transmission string  `json:"transmission"`
	Drivetrain   string  `json:"drivetrain"`
	EngineSize   float64 `json:"engine_size"`
	EngineType   string  `json:"engine_type"`
	HorsePower   uint16  `json:"horse_power"`
	Torque       uint16  `json:"torque"`
}

type CarModsOutputDTO struct {
	Id         uint64 `json:"id"`
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
}
