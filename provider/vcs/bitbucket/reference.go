package bitbucket

// Reference struct
type Reference struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Target *Commit `json:"target"`
	Links  *Links  `json:"links"`
}

/*
{
          "type": "branch",
          "name": "name-of-branch",
          "target": {
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
          },
          "links": {
            "self": {
              "href": "https://api.bitbucket.org/2.0/repositories/user_name/repo_name/refs/branches/master"
            },
            "commits": {
              "href": "https://api.bitbucket.org/2.0/repositories/user_name/repo_name/commits/master"
            },
            "html": {
              "href": "https://bitbucket.org/user_name/repo_name/branch/master"
            }
          }
        }
*/
