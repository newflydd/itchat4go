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
		FatherName:  "1-0|编程",
		Description: "这些群组可以快速帮你结识相关领域志同道合的爱好者，程序猿并不是一个人在战斗！",
		ExampleStr:  `"Golang"或者其索引号"1-1"，我会邀请您加入专业探讨Golang技术的微信群。`,
		ChildrenName: []string{"1-1|Golang", "1-2|Java", "1-3|Python", "1-4|Nodejs", "1-5|Qt",
			"1-6|SQL|数据库",
			"1-7|Angular", "1-8|Vue", "1-9|React", "1-10|jQuery",
			"1-11|Linux", "1-12|Android", "1-13|IOS", "1-14|全栈"}})

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:   "2-0|嵌入式",
		Description:  "这些群组可以快速帮你结识嵌入式方向的爱好者，一起碰撞你们创意的火花吧！",
		ExampleStr:   `"单片机"或者"STM32"或者其索引号"2-1"，我会邀请您加入专业探讨单片机技术的微信群。`,
		ChildrenName: []string{"2-1|STC|STM32|ARM|单片机", "2-2|树莓派|Raspberry|Ardunio|智能硬件"}})

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:   "3-0|互联网",
		Description:  "如果你是下面这些互联网品牌的深度用户，请团结在一起互粉吧！",
		ExampleStr:   `"小米"或者"米粉"或者其索引号"3-1"，我会邀请您加入由全国各地米粉组建的微信群。`,
		ChildrenName: []string{"3-1|小米|米粉", "3-2|知乎|知友", "3-3|微信", "3-4|微博"}})

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:   "4-0|足球",
		Description:  "足球是引人入胜的，那些支持多年的俱乐部让人尤为动情，如果你喜欢他们，加入群组一起讨论吧！",
		ExampleStr:   `"尤文"或者其索引号"4-1"，我会邀请您加入由全国各地尤文蒂尼组建的微信群。`,
		ChildrenName: []string{"4-1|尤文", "4-2|AC", "4-3|国际", "4-4|巴萨", "4-5|拜仁", "4-6|曼联", "4-7|足球赛事"}})

	focusGroupKeywords = append(focusGroupKeywords, FocusGroupKeywords{
		FatherName:   "5-0|同城",
		Description:  "如果你身在这些城市，或许参与进来聊天会发生很多美好的故事...",
		ExampleStr:   `"北京"或者其索引号"5-1"，我会邀请您加入由很多在北京打拼的小伙伴组建的同城微信群`,
		ChildrenName: []string{"5-1|北京", "5-2|上海", "5-3|南京", "5-4|深圳", "5-5|杭州", "5-6|泰州"}})
}

func GetFocusGroupKeywords() []FocusGroupKeywords {
	return focusGroupKeywords
}

/* 获取所有关键词 */
func GetFocusGroupKeywordChildren() []string {
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
		keywordsStr += v.FatherName + "\n"
	}

	return keywordsStr
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
		keywordsStr += str + "\n"
	}

	return focusGroupKeywords[index].Description, focusGroupKeywords[index].ExampleStr, keywordsStr
}
