package chat_repository

import (
	"strconv"
	"strings"
)

const ChatStateBucket string = "chat_state"

type State struct {
	Branch string
	Step   int
}

// NewState format of incoming state: [branch_name].[step]
func NewState(st string) State {
	s := strings.Split(st, ".")
	if len(s) < 2 {
		return State{Branch: s[0], Step: 0}
	}
	step, err := strconv.Atoi(s[1])
	if err != nil {
		return State{Branch: s[0], Step: 0}
	}

	return State{Branch: s[0], Step: step}
}

func (s State) String() string {
	step := strconv.Itoa(s.Step)
	return s.Branch + step
}

type ChatRepository interface {
	Get(userId int64, bucket string) (*State, error)
	Update(userID int64, step string, bucket string) error
	Delete(userID int64, bucket string) error
}
