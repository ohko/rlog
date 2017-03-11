package rlog

import "encoding/json"

func enData(key, value string) []byte {
	d := map[string]string{"k": key, "v": value}
	js, _ := json.Marshal(d)
	return js
}

func deData(data string) (string, string) {
	var d map[string]string
	json.Unmarshal([]byte(data), &d)
	return d["k"], d["v"]
}
