package bitbucket

// Repository struct
type Repository struct {
	Type      string   `json:"type"`
	Links     *Links   `json:"links"`
	UUID      string   `json:"uuid"`
	Project   *Project `json:"project"`
	FullName  string   `json:"full_name"`
	Name      string   `json:"name"`
	Website   string   `json:"website,omitempty"`
	Owner     *Owner   `json:"owner"`
	SCM       string   `json:"scm"`
	IsPrivate bool     `json:"is_private"`
}

/*
{
  "type": "repository",
  "links": {
    "self": {
      "href": "https://api.bitbucket.org/api/2.0/repositories/bitbucket/bitbucket"
    },
    "html": {
      "href": "https://api.bitbucket.org/bitbucket/bitbucket"
    },
    "avatar": {
      "href": "https://api-staging-assetroot.s3.amazonaws.com/c/photos/2014/Aug/01/bitbucket-logo-2629490769-3_avatar.png"
    }
  },
  "uuid": "{673a6070-3421-46c9-9d48-90745f7bfe8e}",
  "project": Project,
  "full_name": "team_name/repo_name",
  "name": "repo_name",
  "website": "https://mywebsite.com/",
  "owner": Owner,
  "scm": "git",
  "is_private": true
},
*/
