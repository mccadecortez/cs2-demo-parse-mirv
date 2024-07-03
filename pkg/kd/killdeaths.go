package kd

type KillDeaths struct {
	// map[steamId32]gametick
	Data map[uint32][]int64 `json:"Inner"`
}

func NewKillDeaths() KillDeaths {
	return KillDeaths{
		Data: make(map[uint32][]int64),
	}
}

func (m *KillDeaths) Add(key uint32, value int64) {
	m.Data[key] = append(m.Data[key], value)
}

func (m *KillDeaths) Get(key uint32) []int64 {
	return m.Data[key]
}
