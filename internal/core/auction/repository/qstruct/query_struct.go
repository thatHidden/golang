package qstruct

type QueryParams struct {
	UserID       uint64
	IsNoReserve  bool
	IsEnded      bool
	IsEndingSoon bool
	MinYear      uint16
	MaxYear      uint16
	Brand        string
	Model        string
	Generation   string
	BodyStyle    string
}
