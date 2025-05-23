package minesweeper

const rowSeparator = byte(0b1111)

func encodeCell(c cell) byte {
	var encoded = byte(0b1000)
	if c.isMine {
		encoded |= 1 << 0 // Set bit 0 for mine
	}
	if c.isFlagged {
		encoded |= 1 << 1 // Set bit 1 for flagged
	}
	if c.isRevealed {
		encoded |= 1 << 2 // Set bit 2 for revealed
	}
	return encoded
	//return byte(0b1001)
}

func decodeCell(encoded byte) cell {
	return cell{
		isMine:     (encoded & (1 << 0)) != 0,
		isFlagged:  (encoded & (1 << 1)) != 0,
		isRevealed: (encoded & (1 << 2)) != 0,
	}
}

func encodeTable(t [][]cell) []byte {
	columns := len(t)
	rows := len(t[0])

	totalLengthBytes := (columns*rows + columns + 1) / 2
	encoded := make([]byte, totalLengthBytes)

	chunkSize := 4
	chunkProgress := 0

	for _, row := range t {
		for _, c := range row {
			byteIndex := (chunkProgress) / 2
			e := encodeCell(c)
			if chunkProgress%2 == 0 {
				encoded[byteIndex] = e
			} else {
				encoded[byteIndex] |= e << chunkSize
			}
			chunkProgress++
		}
		byteIndex := (chunkProgress) / 2
		if chunkProgress%2 == 0 {
			encoded[byteIndex] = rowSeparator
		} else {
			encoded[byteIndex] |= rowSeparator << chunkSize
		}
		chunkProgress++
	}

	return encoded
}

func decodeTable(encoded []byte) [][]cell {
	t := make([][]cell, 1)

	lastStepNewRow := true

	for _, b := range encoded {
		first4Bits := b & 0b1111
		if first4Bits == rowSeparator {
			if lastStepNewRow {
				break
			}
			t = append(t, make([]cell, 0))
			lastStepNewRow = true
		} else {
			t[len(t)-1] = append(t[len(t)-1], decodeCell(first4Bits))
			lastStepNewRow = false
		}

		second4Bits := (b >> 4) & 0b1111
		if second4Bits == rowSeparator {
			if lastStepNewRow {
				break
			}
			t = append(t, make([]cell, 0))
			lastStepNewRow = true
		} else {
			t[len(t)-1] = append(t[len(t)-1], decodeCell(second4Bits))
			lastStepNewRow = false
		}
	}
	// If the last row is empty, remove it
	if len(t[len(t)-1]) == 0 {
		t = t[:len(t)-1]
	}
	return t
}
