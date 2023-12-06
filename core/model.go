//极简博客

package core

//菜单项结构体定义
type MenuItem struct {
	//路径
	Path      string `json:"path"`
	//名称
	Name      string `json:"name"`
	//前端解析模块
	Component string `json:"component"`
}

//文章结构体定义
type Article struct {
	//标题
	Title    string `json:"title"`
	//副标题
	SubTitle string `json:"subTitle"`
	//图表
	Avatar   string `json:"avatar"`
	//所属组
	Group    string `json:"group"`
}
