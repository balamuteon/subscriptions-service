package domain

type Subscription struct {
	ID          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     *string
}
