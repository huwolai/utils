package page

type Page struct  {

	PageSize uint64 `json:"page_size"`
	PageIndex uint64 `json:"page_index"`

	Data interface{} `json:"data"`

}

func NewPage(pageIndex uint64,pageSize uint64,data interface{}) *Page {

	return &Page{PageIndex:pageIndex,PageSize:pageSize,Data:data}
}
