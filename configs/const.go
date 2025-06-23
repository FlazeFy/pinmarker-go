package configs

var ResponseMessages = map[string]string{
	"post":        "created",
	"put":         "updated",
	"hard delete": "permanentally deleted",
	"soft delete": "deleted",
	"recover":     "recovered",
	"get":         "fetched",
	"login":       "login",
	"sign out":    "signed out",
}

var TrackTypes = []string{"live", "share-loc"}
var AppsSources = []string{"pinmarker", "mi-fik", "myride", "kumande"}

// Doc Name
var TrackDoc = "tracks"
