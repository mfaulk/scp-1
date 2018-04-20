package scp

type Msg interface {
	BN() int // returns ballot.counter in PREPARE and COMMIT messages
	Less(Msg) bool
}

// NomMsg is the payload of a nomination protocol message.
type NomMsg struct {
	X, Y ValueSet
}

func (nm *NomMsg) BN() int { return 0 }

func (nm *NomMsg) Less(other Msg) bool {
	if other, ok := other.(*NomMsg); ok {
		if len(nm.Y) < len(other.Y) {
			return true
		}
		if len(other.Y) < len(nm.Y) {
			return false
		}
		return len(nm.X) < len(other.X)
	}
	return true
}

type PrepMsg struct {
	B, P, PP Ballot
	HN, CN   int
}

func (pm *PrepMsg) BN() int { return pm.B.N }

func (pm *PrepMsg) Less(other Msg) bool {
	switch other := other.(type) {
	case *NomMsg:
		return false
	case *PrepMsg:
		if pm.B.Less(other.B) {
			return true
		}
		if other.B.Less(pm.B) {
			return false
		}
		if pm.P.Less(other.P) {
			return true
		}
		if other.P.Less(pm.P) {
			return false
		}
		if pm.PP.Less(other.PP) {
			return true
		}
		if other.PP.Less(pm.PP) {
			return false
		}
		return pm.HN < other.HN
	}
	return true
}

type CommitMsg struct {
	B          Ballot
	PN, HN, CN int
}

func (cm *CommitMsg) BN() int { return cm.B.N }

func (cm *CommitMsg) Less(other Msg) bool {
	switch other := other.(type) {
	case *NomMsg:
		return false
	case *PrepMsg:
		return false
	case *CommitMsg:
		if cm.B.Less(other.B) {
			return true
		}
		if other.B.Less(cm.B) {
			return false
		}
		if cm.PN < other.PN {
			return true
		}
		if other.PN < cm.PN {
			return false
		}
		return cm.HN < other.HN
	}
	return true
}

type ExtMsg struct {
	C  Ballot
	HN int
}

func (em *ExtMsg) BN() int { return 0 }

func (em *ExtMsg) Less(other Msg) bool {
	if other, ok := other.(*ExtMsg); ok {
		return em.HN < other.HN
	}
	return false
}
