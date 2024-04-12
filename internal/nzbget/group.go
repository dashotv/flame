package nzbget

type Group struct {
	ID                 int    `json:"nzbid"` // 4
	RemainingSizeMB    int    // 3497
	PausedSizeMB       int    // 3497
	RemainingFileCount int    // 73
	RemainingParCount  int    // 9
	MinPriority        int    // 0
	MaxPriority        int    // 0
	ActiveDownloads    int    // 0
	Status             string // PAUSED
	NZBName            string // Brave.Are.the.Fallen.2020.1080p.AMZN.WEB-DL.DDP2.0.H.264-ExREN,
	NZBNicename        string // Brave.Are.the.Fallen.2020.1080p.AMZN.WEB-DL.DDP2.0.H.264-ExREN,
	Kind               string // NZB
	URL                string // ,
	NZBFilename        string // Brave.Are.the.Fallen.2020.1080p.AMZN.WEB-DL.DDP2.0.H.264-ExREN,
	DestDir            string // /data/intermediate/Brave.Are.the.Fallen.2020.1080p.AMZN.WEB-DL.DDP2.0.H.264-ExREN.#4,
	FinalDir           string // ,
	Category           string // ,
	ParStatus          string // NONE
	ExParStatus        string // NONE
	UnpackStatus       string // NONE
	MoveStatus         string // NONE
	ScriptStatus       string // NONE
	DeleteStatus       string // NONE
	MarkStatus         string // NONE
	UrlStatus          string // NONE
	FileSizeMB         int    // 3651
	FileCount          int    // 77
	MinPostTime        int    // 1586073677
	MaxPostTime        int    // 1586073793
	TotalArticles      int    // 4992
	SuccessArticles    int    // 212
	FailedArticles     int    // 0
	Health             int    // 1000
	CriticalHealth     int    // 898
	DupeKey            string //
	DupeScore          int    // 0
	DupeMode           string // SCORE
	Deleted            bool
	DownloadedSizeMB   int // 235
	DownloadTimeSec    int // 44
	PostTotalTimeSec   int // 0
	ParTimeSec         int // 0
	RepairTimeSec      int // 0
	UnpackTimeSec      int // 0
	MessageCount       int // 95
	ExtraParBlocks     int // 0
	Parameters         []Parameter
	ScriptStatuses     []ScriptStatus
	ServerStats        []ServerStat
	PostInfoText       string // NONE
	PostStageProgress  int    // 9193728
	PostStageTimeSec   int    // 0
	Log                []Log
	//FirstID            int    // 4
	//LastID             int    // 4
	//RemainingSizeLo    int    // 3666882216
	//RemainingSizeHi    int    // 0
	//PausedSizeLo       int    // 3666882216
	//PausedSizeHi       int    // 0
	//FileSizeLo         int    // 3829038352
	//FileSizeHi         int    // 0
	//DownloadedSizeLo   int // 247289836
	//DownloadedSizeHi   int // 0
}

type ScriptStatus struct {
	Name   string
	Status string
}

type Log struct {
}

type GroupResponse struct {
	*Response
	Result []Group
}
