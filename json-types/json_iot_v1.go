package json_types


type IotMsgV1 struct {
	Type string `json:"type"`
	Cls string `json:"cls"`
	Subcls string `json:"subcls"`
	Def struct {
		Value interface{} `json:"value"`
		Unit string `json:"unit"`
	} `json:"def"`
	Props map[string]interface{} `json:"props"`
	Ctime string `json:"ctime"`
	UUID string `json:"uuid"`
	Ver float32 `json:"ver"`
	Transp string `json:"transp"`
	Corid string `json:"corid"`
}
