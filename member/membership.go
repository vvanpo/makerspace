package member

import (
	"fmt"
	"strings"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
	"log"
)

func (m *Member) Request_membership(rate string) error {
	plan, ok := m.Plans[rate]
	if !ok {
		return fmt.Errorf("Invalid membership rate")
	}
	if rate == "membership-student" && m.Student == nil {
		return fmt.Errorf("Non-students cannot apply for a student membership rate")
	}
	if ms := m.Get_membership(); ms != nil {
		if ms.Plan.ID == rate {
			return fmt.Errorf("Cannot duplicate existing membership")
		}
		if ms.Plan.Amount < plan.Amount {
			params := &stripe.SubParams{Plan: rate}
			return m.Update_membership(params)
		}
	}
	return m.Request_subscription(plan.ID)
}

func (m *Member) Get_pending_membership() *Pending_subscription {
	for _, p := range m.Get_pending_subscriptions() {
		if strings.HasPrefix(p.Plan_id, "membership") {
			return p
		}
	}
	return nil
}

func (m *Member) Get_membership() *stripe.Sub {
	for _, s := range m.Active_subscriptions() {
		if strings.HasPrefix(s.Plan.ID, "membership") {
			return s
		}
	}
	return nil
}

func (m *Member) Membership_plan() string {
	if ms := m.Get_membership(); ms != nil {
		return ms.Plan.ID
	}
	return ""
}

func (m *Member) Membership_id() string {
	if ms := m.Get_membership(); ms != nil {
		return ms.ID
	}
	return ""
}

func (m *Member) Update_membership(params *stripe.SubParams) error {
	ms := m.Get_membership()
	if ms == nil {
		return fmt.Errorf("No existing membership")
	}
	_, err := sub.Update(ms.ID, params)
	return err
}

func (m *Member) Cancel_membership() {
	ms := m.Get_membership()
	if m != nil {
		_, err := sub.Cancel(ms.ID, nil)
		if err != nil {
			log.Println(err)
		}
	}
	if m.Talk_user() != nil {
		m.Talk_user().Remove_from_group("Members")
	}
}