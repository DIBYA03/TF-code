package business

type kycQueryStatus string

var kycQuery kycQueryStatus

const (
	kycQueryStatusNotStarted = kycQueryStatus("notStarted")
	kycQueryStatusDeclined   = kycQueryStatus("declined")
	kycQueryStatusApproved   = kycQueryStatus("approved")
	kycQueryStatusInReview   = kycQueryStatus("review")
)

var kycQueries = map[kycQueryStatus]kycQueryStatus{
	kycQueryStatusNotStarted: kycQueryStatusNotStarted,
	kycQueryStatusDeclined:   kycQueryStatusDeclined,
	kycQueryStatusApproved:   kycQueryStatusApproved,
	kycQueryStatusInReview:   kycQueryStatusInReview,
}

func (q kycQueryStatus) Valid() bool {
	_, ok := kycQueries[q]
	return ok
}

func (kycQueryStatus) new(v string) (kycQueryStatus, bool) {
	q := kycQueryStatus(v)
	return q, q.Valid()
}

func (q kycQueryStatus) String() string {
	return string(q)
}
