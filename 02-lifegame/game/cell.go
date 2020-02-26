package game

// Cell :１マスのセルを表現する
type Cell interface {
	// 次の世代を作成する
	Next(topLeft Cell, topCenter Cell, topRight Cell,
		middleLeft Cell, middleRight Cell,
		bottomLeft Cell, bottomCenter Cell, bottomRight Cell) Cell

	// 生存状態
	IsAlive() bool
}

// AliveCell :生きているセルを表現する
type AliveCell struct {
}

// IsAlive : 生存状態
func (c AliveCell) IsAlive() bool {
	return true
}

// Next : 次の世代のセルを作成する
func (c AliveCell) Next(topLeft Cell, topCenter Cell, topRight Cell,
	middleLeft Cell, middleRight Cell,
	bottomLeft Cell, bottomCenter Cell, bottomRight Cell) Cell {
	var count = 0
	var cells = [...]Cell{topLeft, topCenter, topRight, middleLeft, middleRight, bottomLeft, bottomCenter, bottomRight}
	for i := 0; i < len(cells); i++ {
		if cells[i].IsAlive() {
			count++
		}
	}

	if count == 2 || count == 3 {
		return AliveCell{}
	}
	return DeadCell{}
}

// DeadCell :死んでいるセルを表現する
type DeadCell struct {
}

// IsAlive : 生存状態
func (c DeadCell) IsAlive() bool {
	return false
}

// Next : 次の世代のセルを作成する
func (c DeadCell) Next(topLeft Cell, topCenter Cell, topRight Cell,
	middleLeft Cell, middleRight Cell,
	bottomLeft Cell, bottomCenter Cell, bottomRight Cell) Cell {
	var count = 0
	var cells = [...]Cell{topLeft, topCenter, topRight, middleLeft, middleRight, bottomLeft, bottomCenter, bottomRight}
	for i := 0; i < len(cells); i++ {
		if cells[i].IsAlive() {
			count++
		}
	}

	if count == 3 {
		return AliveCell{}
	}
	return DeadCell{}
}

// EdgeCell :画面端のセルを表現する
type EdgeCell struct {
}

// IsAlive : 生存状態
func (c EdgeCell) IsAlive() bool {
	return false
}

// Next : 次の世代のセルを作成する
func (c EdgeCell) Next(topLeft Cell, topCenter Cell, topRight Cell,
	middleLeft Cell, middleRight Cell,
	bottomLeft Cell, bottomCenter Cell, bottomRight Cell) Cell {
	return EdgeCell{}
}
