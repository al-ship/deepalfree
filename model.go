package deepalfree

type Response struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}
type ApkInfo struct {
	ApkName  string `json:"apk_name"`
	FileSize string `json:"file_size"`
	HashCode string `json:"hash_code"`
	HashType string `json:"hash_type"`
	URL      string `json:"url"`
}
type PayInfo struct {
	Discount      float64 `json:"discount"`
	OrderSource   string  `json:"order_source"`
	OriginalPrice float64 `json:"original_price"`
	PayWays       []any   `json:"pay_ways"`
	Price         float64 `json:"price"`
}
type Statics struct {
	Downloads  string `json:"downloads"`
	Installs   string `json:"installs"`
	Uninstalls string `json:"uninstalls"`
	Updates    string `json:"updates"`
}
type Tags struct {
	TagCode  string `json:"tag_code"`
	TypeCode string `json:"type_code"`
}
type App struct {
	ApkInfo         ApkInfo `json:"apk_info"`
	AppID           string  `json:"app_id"`
	Icon            string  `json:"icon"`
	Name            string  `json:"name"`
	PackageName     string  `json:"package_name"`
	PayInfo         PayInfo `json:"pay_info"`
	RestrictedState int     `json:"restricted_state"`
	Slogan          string  `json:"slogan"`
	Statics         Statics `json:"statics"`
	Tags            []Tags  `json:"tags"`
	Tid             string  `json:"tid"`
	Type            string  `json:"type"`
	Uninstall       bool    `json:"uninstall"`
	Version         string  `json:"version"`
	VersionID       string  `json:"version_id"`
	VersionNumber   int     `json:"version_number"`
}
type Pagination struct {
	Current  int `json:"current"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
type Data struct {
	List       []App      `json:"list"`
	Pagination Pagination `json:"pagination"`
}
