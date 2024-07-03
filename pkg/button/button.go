package button

var (
	BUTTON_ATTACK    uint64 = (1 << 0)  // Fire weapon
	BUTTON_JUMP      uint64 = (1 << 1)  // Jump
	BUTTON_DUCK      uint64 = (1 << 2)  // Crouch
	BUTTON_FORWARD   uint64 = (1 << 3)  // Walk forward
	BUTTON_BACK      uint64 = (1 << 4)  // Walk backwards
	BUTTON_USE       uint64 = (1 << 5)  // Use = (Defuse bomb, etc...)
	BUTTON_CANCEL    uint64 = (1 << 6)  // ??
	BUTTON_LOOKLEFT  uint64 = (1 << 7)  // +left (look)
	BUTTON_LOOKRIGHT uint64 = (1 << 8)  // +right (look)
	BUTTON_MOVELEFT  uint64 = (1 << 9)  // Alias? = (not sure)
	BUTTON_MOVERIGHT uint64 = (1 << 10) // Alias? = (not sure)
	BUTTON_ATTACK2   uint64 = (1 << 11) // Secondary fire = (Revolver, Glock change fire mode, Famas change fire mode) = (not sure)
	BUTTON_RUN       uint64 = (1 << 12)
	BUTTON_RELOAD    uint64 = (1 << 13) // Reload weapon
	BUTTON_ALT1      uint64 = (1 << 14) // +alt2
	BUTTON_ALT2      uint64 = (1 << 15) // +alt1
	BUTTON_WALK      uint64 = (1 << 16) // Shift
	BUTTON_SPEED     uint64 = (1 << 17) // Sprint in other source games
	BUTTON_SCORE     uint64 = (1 << 18)
	BUTTON_ZOOM      uint64 = (1 << 19) // Zoom weapon = (not sure)
	BUTTON_WEAPON1   uint64 = (1 << 20)
	BUTTON_WEAPON2   uint64 = (1 << 21)
	BUTTON_BULLRUSH  uint64 = (1 << 22)
	BUTTON_INSPECT   uint64 = (1 << 35)
	// guesses:
	// 33 == quickswitch
)

type Buttons struct {
	Attack    bool `json:"attack"`
	Jump      bool `json:"jump"`
	Duck      bool `json:"duck"`
	Forward   bool `json:"forward"`
	Back      bool `json:"back"`
	Use       bool `json:"use"`
	Cancel    bool `json:"cancel"`
	LookLeft  bool `json:"look_left"`
	LookRight bool `json:"look_right"`
	MoveLeft  bool `json:"move_left"`
	MoveRight bool `json:"move_right"`
	Attack2   bool `json:"attack2"`
	Run       bool `json:"run"`
	Reload    bool `json:"reload"`
	Alt1      bool `json:"alt1"`
	Alt2      bool `json:"alt2"`
	Walk      bool `json:"walk"`
	Speed     bool `json:"speed"`
	Score     bool `json:"score"`
	Zoom      bool `json:"zoom"`
	Weapon1   bool `json:"weapon1"`
	Weapon2   bool `json:"weapon2"`
	Bullrush  bool `json:"bullrush"`
	Inspect   bool `json:"inspect"`
}

func ParseButtonMask(mask uint64) *Buttons {
	attack := (mask & BUTTON_ATTACK) != 0
	jump := (mask & BUTTON_JUMP) != 0
	duck := (mask & BUTTON_DUCK) != 0
	forward := (mask & BUTTON_FORWARD) != 0
	back := (mask & BUTTON_BACK) != 0
	use := (mask & BUTTON_USE) != 0
	cancel := (mask & BUTTON_CANCEL) != 0
	lookleft := (mask & BUTTON_LOOKLEFT) != 0
	lookright := (mask & BUTTON_LOOKRIGHT) != 0
	moveleft := (mask & BUTTON_MOVELEFT) != 0
	moveright := (mask & BUTTON_MOVERIGHT) != 0
	attack2 := (mask & BUTTON_ATTACK2) != 0
	run := (mask & BUTTON_RUN) != 0
	reload := (mask & BUTTON_RELOAD) != 0
	alt1 := (mask & BUTTON_ALT1) != 0
	alt2 := (mask & BUTTON_ALT2) != 0
	walk := (mask & BUTTON_WALK) != 0
	speed := (mask & BUTTON_SPEED) != 0
	score := (mask & BUTTON_SCORE) != 0
	zoom := (mask & BUTTON_ZOOM) != 0
	weapon1 := (mask & BUTTON_WEAPON1) != 0
	weapon2 := (mask & BUTTON_WEAPON2) != 0
	bullrush := (mask & BUTTON_BULLRUSH) != 0
	inspect := (mask & BUTTON_INSPECT) != 0

	return &Buttons{
		attack,
		jump,
		duck,
		forward,
		back,
		use,
		cancel,
		lookleft,
		lookright,
		moveleft,
		moveright,
		attack2,
		run,
		reload,
		alt1,
		alt2,
		walk,
		speed,
		score,
		zoom,
		weapon1,
		weapon2,
		bullrush,
		inspect,
	}
}
