package entities

import (
	"sort"
	"time"
)

type AuthenticatorAssuranceLevel int

const (
	Aal1 AuthenticatorAssuranceLevel = iota
	Aal2
	Aal3
)

func (aal AuthenticatorAssuranceLevel) String() string {
	switch aal {
	case Aal1:
		return "aal1"
	case Aal2:
		return "aal2"
	case Aal3:
		return "aal3"
	default:
		return ""
	}
}

type SessionValidityReason = int

const (
	SessionValid SessionValidityReason = iota << 1
	SessionPastNotAfter
	SessionPastTimebox
	SessionTimedOut
)

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
	CreatedAt   time.Time      `xorm:"created TIMESTAMPZ created_at"`
	DeletedAt   *time.Time     `xorm:"TIMESTAMPZ deleted_at"`
	LastUsedAt  *time.Time     `xorm:"TIMESTAMPZ last_used_at"`
	NotAfter    *time.Time     `xorm:"null TIMESTAMPZ not_after"`
	RefreshedAt *time.Time     `xorm:"null TIMESTAMPZ refreshed_at"`
	ID          string         `xorm:"not null pk VARCHAR(36) id"`
	UserID      string         `xorm:"not null VARCHAR(36) user_id"`
	AAL         string         `xorm:"not null TEXT aal"`
	IP          string         `xorm:"not null TEXT default '' ip"`
	UserAgent   string         `xorm:"not null TEXT default '' user_agent"`
	FactorID    string         `xorm:"not null VARCHAR(36) default '' factor_id"`
	AMRClaims   []*MFAAMRClaim `xorm:"-"`
}

func (Session) TableName() string {
	return "sessions"
}

func (s *Session) CalculateAALAndAMR(
	user *User,
) (aal AuthenticatorAssuranceLevel, amr []AMREntry, err error) {
	amr, aal = []AMREntry{}, Aal1
	for _, claim := range s.AMRClaims {
		if claim.AuthenticationMethod == TOTPSignIn.String() {
			aal = Aal2
		}
		amr = append(
			amr,
			AMREntry{Method: claim.GetAuthenticationMethod(), Timestamp: claim.UpdatedAt.Unix()},
		)
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

func (s *Session) CheckValidity(
	now time.Time,
	refreshTokenTime *time.Time,
	timebox, inactivityTimeout *time.Duration,
) SessionValidityReason {
	if s.NotAfter != nil && now.After(*s.NotAfter) {
		return SessionPastNotAfter
	}

	if timebox != nil && *timebox != 0 && now.After(s.CreatedAt.Add(*timebox)) {
		return SessionPastTimebox
	}

	if inactivityTimeout != nil && *inactivityTimeout != 0 &&
		now.After(s.LastRefreshedAt(refreshTokenTime).Add(*inactivityTimeout)) {
		return SessionTimedOut
	}

	return SessionValid
}

func (s *Session) LastRefreshedAt(refreshTokenTime *time.Time) time.Time {
	refreshedAt := s.RefreshedAt

	if refreshedAt == nil || refreshedAt.IsZero() {
		if refreshTokenTime != nil {
			rtt := *refreshTokenTime

			if rtt.IsZero() {
				return s.CreatedAt
			} else if rtt.After(s.CreatedAt) {
				return rtt
			}
		}

		return s.CreatedAt
	}

	return *refreshedAt
}
