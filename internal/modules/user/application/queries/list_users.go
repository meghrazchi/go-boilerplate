package queries

type ListUsersQuery struct {
	Page   int
	Limit  int
	Search string
	Sort   string
	Order  string
}
