package content

type Content0200 struct {
	Alarm     []int       `json:"alarm"`
	Dir       int         `json:"dir"`
	GPSTime   string      `json:"gpsTime"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	Mileage   int         `json:"mileage"`
	OnlineNo  string      `json:"onlineNo"`
	RouteCode string      `json:"routeCode"`
	Speed     int         `json:"speed"`
	StaIndex  int         `json:"staIndex"`
	State     map[int]int `json:"state"`
	UpDown    int         `json:"upDown"`
}

type Content0DE0 struct {
	PlanDate string `json:"planDate"`
	PlanList []int  `json:"planList"`
	State    string `json:"state"`
}

type Content0202 struct {
	ID uint64 `json:"id" gorm:"primaryKey;autoIncrement"`

	OnlineNo  string        `json:"onlineNo"`
	Longitude float64       `json:"longitude"`
	Latitude  float64       `json:"latitude"`
	Speed     int           `json:"speed"`
	Dir       int           `json:"dir"`
	GPSTime   string        `json:"gpsTime"`
	CanItem   map[byte]byte `json:"canItem"`
}
