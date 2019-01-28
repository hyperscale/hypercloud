package bitbucket

// Owner struct
type Owner struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
	Links       *Links `json:"links"`
}

/*
{
  "type": "user",
  "username": "emmap1",
  "display_name": "Emma",
  "uuid": "{a54f16da-24e9-4d7f-a3a7-b1ba2cd98aa3}",
  "links": {
    "self": {
      "href": "https://api.bitbucket.org/api/2.0/users/emmap1"
    },
    "html": {
      "href": "https://api.bitbucket.org/emmap1"
    },
    "avatar": {
      "href": "https://bitbucket-api-assetroot.s3.amazonaws.com/c/photos/2015/Feb/26/3613917261-0-emmap1-avatar_avatar.png"
    }
  }
},
*/
