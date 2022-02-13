package vars

var (
	ChangeGenderEnabled = true

	WhiteDomainEnabled = false
	WhiteDomain        = ""
	WhiteEmailEnabled  = false
)

var (
	VerifyProfileQuestionId int64
	VerifyOptionId          int64

	MatchingQuestionId         int64
	MatchingVerifiedOptionId   int64
	MatchingUnverifiedOptionId int64
	MatchingAnyOptionId        int64

	FemaleMatchRateLimit MatchRateLimit
	MaleMatchRateLimit   MatchRateLimit
)

type MatchRateLimit struct {
	OptionId            int64
	RateLimitEnabled    bool
	RateLimitUnit       string
	RateLimitUnitPeriod int64
	MatchPerRate        int64
}
