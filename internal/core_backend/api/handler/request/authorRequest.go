package request

type AuthorRequest struct {
	AuthorID string `uri:"author_id" binding:"mongodb"`
}
