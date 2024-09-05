package loader

type TagInfo struct {
	AssistHashIndexDict map[string]*AssistIndex
	AssistTreeIndexDict map[string]*AssistIndex
	TagDict             map[string][]string
	DescRow             []string
	FlagRow             []string
	TypeRow             []string
	KeyRow              []string
	TagRow              []string
}

func (t *TagInfo) init(rows [][]string) {
	t.AssistHashIndexDict = make(map[string]*AssistIndex)
	t.AssistTreeIndexDict = make(map[string]*AssistIndex)
	t.TagDict = make(map[string][]string)
	t.DescRow = make([]string, len(rows[0]))
	t.FlagRow = make([]string, len(rows[1]))
	t.TypeRow = make([]string, len(rows[2]))
	t.KeyRow = make([]string, len(rows[3]))
	t.TagRow = make([]string, len(rows[3]))
}
