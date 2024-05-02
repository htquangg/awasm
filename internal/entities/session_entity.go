package entities

import (
	"sort"
	"time"
)

type AuthenticatorAssuranceLevel int

const (
	AAL1 AuthenticatorAssuranceLevel = iota
	AAL2
	AAL3
)

func (aal AuthenticatorAssuranceLevel) String() string {
	switch aal {
	case AAL1:
		return "aal1"
	case AAL2:
		return "aal2"
	case AAL3:
		return "aal3"
	default:
		return ""
	}
}

type AMREntry struct {
	Method    string `json:"method"`
	Provider  string `json:"provider,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

type sortAMREntries struct {
	Array []AMREntry
}

func (s sortAMREntries) Len() int {
	return len(s.Array)
}

func (s sortAMREntries) Less(i, j int) bool {
	return s.Array[i].Timestamp < s.Array[j].Timestamp
}

func (s sortAMREntries) Swap(i, j int) {
	s.Array[j], s.Array[i] = s.Array[i], s.Array[j]
}

type Session struct {
	CreatedAt  time.Time      `xorm:"created TIMESTAMPZ created_at"`
	DeletedAt  *time.Time     `xorm:"TIMESTAMPZ deleted_at"`
	LastUsedAt *time.Time     `xorm:"TIMESTAMPZ last_used_at"`
	NotAfter   *time.Time     `xorm:"not null TIMESTAMPZ not_after"`
	ID         string         `xorm:"not null pk VARCHAR(36) id"`
	UserID     string         `xorm:"not null VARCHAR(36) user_id"`
	AAL        string         `xorm:"not null TEXT aal"`
	IP         string         `xorm:"not null TEXT default '' ip"`
	UserAgent  string         `xorm:"not null TEXT default '' user_agent"`
	FactorID   string         `xorm:"not null VARCHAR(36) default '' factor_id"`
	AMRClaims  []*MFAAMRClaim `xorm:"-"`
}

func (Session) TableName() string {
	return "sessions"
}

func (s *Session) CalculateAALAndAMR(user *User) (aal AuthenticatorAssuranceLevel, amr []AMREntry, err error) {
	amr, aal = []AMREntry{}, AAL1
	for _, claim := range s.AMRClaims {
		if claim.AuthenticationMethod == TOTPSignIn.String() {
			aal = AAL2
		}
		amr = append(amr, AMREntry{Method: claim.GetAuthenticationMethod(), Timestamp: claim.UpdatedAt.Unix()})
	}

	// makes sure that the AMR claims are always ordered most-recent first

	// sort in ascending order
	sort.Sort(sortAMREntries{
		Array: amr,
	})

	// now reverse for descending order
	_ = sort.Reverse(sortAMREntries{
		Array: amr,
	})

	return aal, amr, nil
}
