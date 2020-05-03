package nzbget

type Status struct {
	RemainingSizeMB     int  // 0
	ForcedSizeMB        int  // 0
	DownloadedSizeMB    int  // 0
	MonthSizeMB         int  // 0
	DaySizeMB           int  // 0
	ArticleCacheMB      int  // 0
	DownloadRate        int  // 0
	AverageDownloadRate int  // 0
	DownloadLimit       int  // 0
	UpTimeSec           int  // 2281
	DownloadTimeSec     int  // 0
	ServerPaused        bool // false
	DownloadPaused      bool // false
	Download2Paused     bool // false
	ServerStandBy       bool // true
	PostPaused          bool // false
	ScanPaused          bool // false
	QuotaReached        bool // false
	FreeDiskSpaceMB     int  // 134539
	ServerTime          int  // 1586063906
	ResumeTime          int  // 0
	FeedActive          bool // false
	QueueScriptCount    int  // 0
	NewsServers         []NewsServer
	//RemainingSizeLo     int  // 0
	//RemainingSizeHi     int  // 0
	//ForcedSizeLo        int  // 0
	//ForcedSizeHi        int  // 0
	//DownloadedSizeLo    int  // 0
	//DownloadedSizeHi    int  // 0
	//MonthSizeLo         int  // 0
	//MonthSizeHi         int  // 0
	//DaySizeLo           int  // 0
	//DaySizeHi           int  // 0
	//ArticleCacheLo      int  // 0
	//ArticleCacheHi      int  // 0
	//ThreadCount         int  // 7
	//ParJobCount         int  // 0
	//PostJobCount        int  // 0
	//UrlCount            int  // 0
	//FreeDiskSpaceLo     int  // 3635539968
	//FreeDiskSpaceHi     int  // 32
}

type NewsServer struct {
	ID     int
	Active bool
}

type StatusResponse struct {
	*Response
	Result *Status
}
