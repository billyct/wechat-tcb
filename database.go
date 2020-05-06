package tcb

const (
	// 请求地址：「新增集合」
	urlDatabaseCollectionAdd = "https://api.weixin.qq.com/tcb/databasecollectionadd"
	// 请求地址：「数据库插入记录」
	urlDatabaseAdd           = "https://api.weixin.qq.com/tcb/databaseadd"
)

// 新增集合请求
type reqDatabaseCollectionAdd struct {
	Env            string `json:"env,omitempty"`
	CollectionName string `json:"collection_name,omitempty"`
}

// 新增集合返回
type resDatabaseCollectionAdd struct {
	ResError
}

// 新增集合
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/database/databaseCollectionAdd.html
func (t *Tcb) DatabaseCollectionAdd(name string) error {
	req := &reqDatabaseCollectionAdd{
		Env:            t.envId,
		CollectionName: name,
	}

	data, err := httpPostWithReq(t.url(urlDatabaseCollectionAdd), req)
	if err != nil {
		return err
	}

	res := &resDatabaseCollectionAdd{}

	err = DecodeApiData("DatabaseCollectionAdd", data, res)
	return err
}

// 数据库插入记录请求
type reqDatabaseAdd struct {
	Env   string `json:"env,omitempty"`
	Query string `json:"query,omitempty"`
}

// 数据库插入记录返回
type resDatabaseAdd struct {
	ResError
	IDList []string `json:"id_list"`
}

// 数据库插入记录
// https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/database/databaseAdd.html
func (t *Tcb) DatabaseAdd(query string) (*resDatabaseAdd, error) {
	req := &reqDatabaseAdd{
		Env:   t.envId,
		Query: query,
	}

	data, err := httpPostWithReq(t.url(urlDatabaseAdd), req)
	if err != nil {
		return nil, err
	}

	res := &resDatabaseAdd{}

	err = DecodeApiData("DatabaseAdd", data, res)
	return res, err
}
