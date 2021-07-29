package response

type DictResponseStruct struct {
	BaseData
	Name      string                  `json:"name"`
	Desc      string                  `json:"desc"`
	Status    uint                    `json:"status"`
	Remark    string                  `json:"remark"`
	DictDatas []DictDataResponseStruct `json:"dictDatas"`
}

type DictDataResponseStruct struct {
	BaseData
	Key      string            `json:"key"`
	Val      string            `json:"val"`
	Attr     string            `json:"attr"`
	Addition string            `json:"addition"`
	Sort     uint              `json:"sort"`
	Status   uint              `json:"status"`
	Remark   string            `json:"remark"`
	DictId   uint              `json:"dictId"`
	Dict     DictResponseStruct `json:"dict"`
}
