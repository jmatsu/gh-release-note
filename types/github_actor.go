package types

type GitHubActor struct {
	Id    string `json:"id"`
	IsBot bool   `json:"is_bot"`
	Login string `json:"login"`
	Name  string `json:"name"`
}
