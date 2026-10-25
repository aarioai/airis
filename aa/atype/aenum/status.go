package aenum

import (
	"github.com/aarioai/airis/pkg/types"
)

type Status int8 // the status field must not default to 0, to avoid accidents

const (
	SysRevoked Status = -128 // revoked; deleted by the system (review may have failed)
	Expired    Status = -121 // invalid or expired; the system deletes it later
	NotExists  Status = -114 // data does not exist --> usually to initialize a field; uncommon
	Deleted    Status = -100 // deleted by the user; visible to nobody
	/********************* values <= -100 are content to delete, so a whole partition can be dropped ************************/
	Supervised  Status = -90 // the system detected an anomaly and locked the account; waiting for the anomaly to clear
	Deleting    Status = -89 // deleting; enters Deleted only after the cooling-off period
	Suspend     Status = -80 // paused or temporarily hidden, set by the user
	Failed      Status = -20 // censor failed: review failed, the user may edit, visible only to the user --> whatever was set before, this is the single failure state; set the real status again on edit
	Pending     Status = -10 // under review; the user keeps edit rights and may resubmit
	Pending1    Status = -9  // censor 1 passed --> staged review removes the user's edit rights; a system lock on anomaly also enters this state
	Pending2    Status = -8
	Pending3    Status = -7
	Pending4    Status = -6
	Pending5    Status = -5
	Pending6    Status = -4
	Pending7    Status = -3
	Pending8    Status = -2
	Pending9    Status = -1
	Created     Status = 0   // newly created, under review, but shown; the user keeps edit rights
	Passed      Status = 1   // review passed
	SysTakeover Status = 127 // locked by the system; the user must not edit. Deletion is allowed: unbind with a soft delete, or hard-delete, depending on the case

)

var (
	deletedStr      = types.Itoa(Deleted)
	pendingLimit    = types.Itoa(Pending - 1)
	MeReadableRange = [2]Status{Failed, SysTakeover} // range shown in the user's own list
)

func NewStatus(sts int8) (Status, bool) {
	s := Status(sts)
	ok := s == SysRevoked || s == Expired || s == NotExists || s == Deleted
	ok = ok || s == Supervised || s == Deleting || s == Suspend || s == Failed
	ok = ok || (s >= Pending && s <= Pending9) || s == Created || s == Passed || s == SysTakeover
	return s, ok
}
func (s Status) Int8() int8     { return int8(s) }
func (s Status) String() string { return types.Itoa(s) }
func (s Status) Is(x int8) bool { return s.Int8() == x }
func (s Status) In(a ...Status) bool {
	for _, sts := range a {
		if sts == s {
			return true
		}
	}
	return false
}

// Boundary between shown and hidden
func (s Status) IsOk() bool  { return s >= Created }
func (s Status) NotOk() bool { return s < Created }

// Whether review has passed
func (s Status) IsPassed() bool  { return s >= Passed }
func (s Status) NotPassed() bool { return s < Passed }

// Created is also a pending-review state
func (s Status) IsPending() bool { return s >= Pending && s <= Created }
func (s Status) IsFailed() bool  { return s <= Failed }
func (s Status) IsDiscard() bool { return s <= Deleting }
func (s Status) NotExits() bool  { return s == NotExists }

// After/Before are exclusive; From/Within are inclusive.
// FromPending is generally used to check whether a status is under review or already approved.
func (s Status) FromPending() bool { return s >= Pending }

// Whether the user themself can see it
func (s Status) MeReadable() bool { return s >= MeReadableRange[0] && s <= MeReadableRange[1] }

// Whether the user may edit or delete
func (s Status) Modifiable() bool { return s.In(Failed, Pending, Created, Passed) }

func StsInvalid(k string) string    { return k + "<0" }
func StsPublic(k string) string     { return k + ">-1" }
func StsPassed(k string) string     { return k + ">0" }
func StsNotPassed(k string) string  { return k + "<1" }
func StsNotDiscard(k string) string { return k + ">-90" }

// Publicly visible statuses: under review, created, and approved
func StsVisible(k string) string { return k + ">" + pendingLimit }

// Strictly under review
func StsPending(k string) string {
	return k + ">" + pendingLimit + " AND " + k + "<0"
}

// Under review or created
func StsPendingC(k string) string {
	return k + ">" + pendingLimit + " AND " + k + "<1"
}

// Visible to oneself
func StsVis(k string) string { return k + ">" + deletedStr }
