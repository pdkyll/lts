package config

type Release struct {
	VersionAPI 	string `json:"version_api"`
	UpdateAtAPI	string `json:"update_at_api"`
}

var Version Release = Release{
	"0.1.1",		// version API
	"2022-11-27",	// update at API
}