package filemeta

import "github.com/viant/toolbox"

type Meta struct {
	Paths  map[string]*FolderInfo
	Assets map[string]int
}

type FolderInfo struct {
	Folder     string
	FileCount  int
	LinesCount int
}

func (m *Meta) Add(URL string, lineCount int) {
	key, _ := toolbox.URLSplit(URL)
	info, ok := m.Paths[key]

	if !ok {
		info = &FolderInfo{}
		m.Paths[key] = info
	}
	assetLineCount, ok := m.Assets[URL]
	if ok {
		info.FileCount++
		info.LinesCount -= assetLineCount
	}
	m.Assets[URL] = lineCount
	info.FileCount++
	info.LinesCount += lineCount
}
