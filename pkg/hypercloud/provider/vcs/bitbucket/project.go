package bitbucket

// Project struct
type Project struct {
	Type    string `json:"type"`
	Project string `json:"project"`
	UUID    string `json:"uuid"`
	Links   *Links `json:"links"`
	Key     string `json:"key"`
}

/*
{
  "type": "project",
  "project": "Untitled project",
  "uuid": "{3b7898dc-6891-4225-ae60-24613bb83080}",
  "links": {
    "html": {
      "href": "https://bitbucket.org/account/user/teamawesome/projects/proj"
    },
    "avatar": {
      "href": "https://bitbucket.org/account/user/teamawesome/projects/proj/avatar/32"
    }
  },
  "key": "proj"
},
*/
