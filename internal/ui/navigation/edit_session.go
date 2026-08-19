package navigation

// EditSession represents pending unsaved changes in the workspace. Entering
// edit mode alone does not activate the session; the owning form starts it
// only after an editable field differs from its original snapshot.
type EditSession struct {
	active  bool
	discard func()
}

// Begin marks a page as dirty and records how its current values are restored
// when navigation discards the session.
func (s *EditSession) Begin(discard func()) {
	if s == nil || s.active {
		return
	}

	s.active = true
	s.discard = discard
}

// End closes the current edit session without invoking its discard callback.
func (s *EditSession) End() {
	if s == nil {
		return
	}

	s.active = false
	s.discard = nil
}

// Active reports whether navigation must ask before leaving the current
// page.
func (s *EditSession) Active() bool {
	return s != nil && s.active
}

// Discard restores the page values and closes the session.
func (s *EditSession) Discard() {
	if s == nil || !s.active {
		return
	}

	discard := s.discard
	s.End()
	if discard != nil {
		discard()
	}
}
