package bitbucket

// Link struct
type Link struct {
	Href string `json:"href"`
}

// Links struct
type Links struct {
	Self    *Link `json:"self,omitempty"`
	HTML    *Link `json:"html,omitempty"`
	Diff    *Link `json:"diff,omitempty"`
	Commits *Link `json:"commits,omitempty"`
	Avatar  *Link `json:"avatar,omitempty"`
}

/*
{
	"html": {
		"href": "https://bitbucket.org/user_name/repo_name/branches/compare/c4b2b7914156a878aa7c9da452a09fb50c2091f2..b99ea6dad8f416e57c5ca78c1ccef590600d841b"
	},
	"diff": {
		"href": "https://api.bitbucket.org/2.0/repositories/user_name/repo_name/diff/c4b2b7914156a878aa7c9da452a09fb50c2091f2..b99ea6dad8f416e57c5ca78c1ccef590600d841b"
	},
	"commits": {
		"href": "https://api.bitbucket.org/2.0/repositories/user_name/repo_name/commits?include=c4b2b7914156a878aa7c9da452a09fb50c2091f2&exclude=b99ea6dad8f416e57c5ca78c1ccef590600d841b"
	}
}
*/
