package game

// Map : ゲーム全体を表現する
type Map struct {
	Cells [][]Cell
}

// CreateMap : Mapの作成と初期化
func CreateMap(init [][]bool) Map {
	var lifeMap Map
	for i := 0; i < len(init); i++ {
		var row []Cell
		for j := 0; j < len(init[i]); j++ {
			if init[i][j] {
				row = append(row, AliveCell{})
			} else {
				row = append(row, DeadCell{})
			}
		}
		lifeMap.Cells = append(lifeMap.Cells, row)
	}

	return lifeMap
}

// Next : 次のターンのMapを生成する
func (m Map) Next() Map {
	var result = Map{}

	result.Cells = make([][]Cell, len(m.Cells), len(m.Cells))

	for i := 0; i < len(m.Cells); i++ {
		result.Cells[i] = make([]Cell, len(m.Cells[i]), len(m.Cells[i]))

		for j := 0; j < len(m.Cells[i]); j++ {
			result.Cells[i][j] = m.Cells[i][j].Next(
				m.safeGet(i-1, j-1), m.safeGet(i-1, j), m.safeGet(i-1, j+1),
				m.safeGet(i, j-1), m.safeGet(i, j+1),
				m.safeGet(i+1, j-1), m.safeGet(i+1, j), m.safeGet(i+1, j+1))
		}
	}

	return result
}

// Display : 表示用の文字列を作成する
func (m Map) Display() string {
	var result string
	for i := 0; i < len(m.Cells); i++ {
		var row string
		for j := 0; j < len(m.Cells[i]); j++ {
			if m.Cells[i][j].IsAlive() {
				row += "■"
			} else {
				row += "□"
			}
		}
		result += row + "\n"
	}
	return result
}

func (m Map) safeGet(i int, j int) Cell {
	if i < 0 || len(m.Cells) <= i {
		return EdgeCell{}
	}
	if j < 0 || len(m.Cells[i]) <= j {
		return EdgeCell{}
	}

	return m.Cells[i][j]
}
