package enum

/**
 * 这个文件用来维护焦点微信群的关键字的正则表达式
 * TODO:后期使用数据库维护
 */

type FocusGroupKeywords struct {
	FatherName   string
	Description  string
	ExampleStr   string
	ChildrenName []string
}

var focusGroupKeywords []FocusGroupKeywords

func init() {
	focusGroupKeywords = []FocusGroupKeywords{}

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:  "编程",
		Description: "这些群组可以快速帮你结识相关领域志同道合的爱好者，程序猿并不是一个人在战斗！",
		ExampleStr:  `"::Golang"，我会邀请您加入专业探讨Golang技术的微信群。`,
		ChildrenName: []string{"Golang", "Java", "Python", "Nodejs", "Qt",
			"MySQL", "MongoDB", "PostgreSQL", "Oracle", "MSSQL",
			"Angular", "Vue", "React", "jQuery",
			"Linux", "MacOS", "Android", "IOS", "全栈"}})

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:   "嵌入式",
		Description:  "这些群组可以快速帮你结识嵌入式方向的爱好者，一起碰撞你们创意的火花吧！",
		ExampleStr:   `"::STC"，我会邀请您加入专业探讨STC系列单片机的微信群。`,
		ChildrenName: []string{"STC", "STM32", "ARM", "Raspberry", "Ardunio", "NodeMCU"}})

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:   "互联网",
		Description:  "如果你是下面这些互联网品牌的深度用户，请团结在一起互粉吧！",
		ExampleStr:   `"::小米"或者"::米粉"，我会邀请您加入由全国各地米粉组建的微信群。`,
		ChildrenName: []string{"(小米|米粉)", "(知乎|知友)", "微信", "微博"}})

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:   "足球",
		Description:  "足球是引人入胜的，那些支持多年的俱乐部让人尤为动情，如果你喜欢他们，加入群组一起讨论吧！",
		ExampleStr:   `"::尤文"，我会邀请您加入由全国各地尤文蒂尼组建的微信群。`,
		ChildrenName: []string{"尤文", "AC", "国际", "巴萨", "拜仁", "曼联", "世界杯", "欧冠"}})

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:   "地区",
		Description:  "如果你身在这些城市，或许参与进来聊天会发生很多故事...",
		ExampleStr:   `"::北京"，我会邀请您加入北京同城微信群`,
		ChildrenName: []string{"北京", "上海", "南京", "深圳", "广州", "杭州", "成都", "泰州"}})
}

/* 获取所有关键词 */
func GetFocusGroupKeywords() []string {
	keywords := []string{}

	for _, v := range focusGroupKeywords {
		for _, str := range v.ChildrenName {
			keywords = append(keywords, str)
		}
	}

	return keywords
}

/* 得到关键词父级目录 */
func GetFatherKeywordsStr() string {
	keywordsStr := ""

	for _, v := range focusGroupKeywords {
		keywordsStr += v.FatherName + "，"
	}

	return keywordsStr[:len(keywordsStr)-3]
}

/* 返回分组介绍和所有分组关键词的组装 */
func GetChildKeywordsInfo(fatherName string) (string, string, string) {
	keywordsStr := ""
	var index int

	for index = 0; index < len(focusGroupKeywords); index++ {
		if focusGroupKeywords[index].FatherName == fatherName {
			break
		}
	}

	if index == len(focusGroupKeywords) {
		return "", "", ""
	}

	for _, str := range focusGroupKeywords[index].ChildrenName {
		keywordsStr += str + "，"
	}

	return focusGroupKeywords[index].Description, focusGroupKeywords[index].ExampleStr, keywordsStr[:len(keywordsStr)-3]
}
