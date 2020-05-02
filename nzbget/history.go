package nzbget

type History struct {
	ID                 int `json:"nzbid"`
	Name               string
	RemainingFileCount int
	RetryData          bool
	HistoryTime        int
	Status             string
	Log                []string
	NZBName            string
	NZBNicename        string
	Kind               string
	URL                string
	NZBFilename        string
	DestDir            string
	FinalDir           string
	Category           string
	ParStatus          string
	ExParStatus        string
	UnpackStatus       string
	MoveStatus         string
	ScriptStatus       string
	DeleteStatus       string
	MarkStatus         string
	UrlStatus          string
	FileSizeLo         int
	FileSizeHi         int
	FileSizeMB         int
	FileCount          int
	MinPostTime        int
	MaxPostTime        int
	TotalArticles      int
	SuccesArticles     int
	FailedArticles     int
	Health             int
	CriticalHealth     int
	DupeKey            string
	DupeScore          int
	DupeMode           string
	Deleted            bool
	DownloadedSizeLo   int
	DownloadedSizeHi   int
	DownloadedSizeMB   int
	DownloadTimeSec    int
	PostTotalTimeSec   int
	ParTimeSec         int
	RepairTimeSec      int
	UnpackTimeSec      int
	MessageCount       int
	ExtraParBlocks     int
	Parameters         []Parameter
	ScriptStatuses     []ScriptStatus
	ServerStats        []ServerStat
}

type Parameter struct {
	Name  string
	Value string
}

type ServerStat struct {
	ServerID        int
	SuccessArticles int
	FailedArticles  int
}

type HistoryResponse struct {
	*Response
	Result []History `json:"Result"`
}
