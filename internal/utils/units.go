package utils


func ConvertToMgdl(value float64) float64 {
	return float64(int64(value*18.0182*100)) / 100
}
