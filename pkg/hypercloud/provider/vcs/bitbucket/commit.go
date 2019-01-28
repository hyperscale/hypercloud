package bitbucket

// Commit struct
type Commit struct {
	Hash    string    `json:"hash"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
	Author  *Owner    `json:"author"`
	Date    string    `json:"date,omitempty"`
	Links   *Links    `json:"links"`
	Parents []*Commit `json:"parents,omitempty"`
}

/*
{
            "type": "commit",
            "hash": "709d658dc5b6d6afcd46049c2f332ee3f515a67d",
            "author": User,
            "message": "new commit message\n",
            "date": "2015-06-09T03:34:49+00:00",
            "parents": [
              {
                "type": "commit",
                "hash": "1e65c05c1d5171631d92438a13901ca7dae9618c",
                "links": {
                  "self": {
                    "href": "https://api.bitbucket.org/2.0/repositories/user_name/repo_name/commit/8cbbd65829c7ad834a97841e0defc965718036a0"
                  },
                  "html": {
                    "href": "https://bitbucket.org/user_name/repo_name/commits/8cbbd65829c7ad834a97841e0defc965718036a0"
                  }
                }
              }
            ],
            "links": {
              "self": {
                "href": "https://api.bitbucket.org/2.0/repositories/user_name/repo_name/commit/c4b2b7914156a878aa7c9da452a09fb50c2091f2"
              },
              "html": {
                "href": "https://bitbucket.org/user_name/repo_name/commits/c4b2b7914156a878aa7c9da452a09fb50c2091f2"
              }
            }
          }
*/

/*
{
	"hash": "03f4a7270240708834de475bcf21532d6134777e",
	"type": "commit",
	"message": "commit message\n",
	"author": User,
	"links": {
		"self": {
			"href": "https://api.bitbucket.org/2.0/repositories/user/repo/commit/03f4a7270240708834de475bcf21532d6134777e"
		},
		"html": {
			"href": "https://bitbucket.org/user/repo/commits/03f4a7270240708834de475bcf21532d6134777e"
		}
	}
}
*/
