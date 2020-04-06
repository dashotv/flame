package nzbget

type VersionResponse struct {
	*Response
	Version string `json:"result"`
}
