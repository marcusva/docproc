package common

import (
// "bytes"
// "encoding/gob"
)

// CopyMap performs a deep copy of an arbitrary map using the encoding/gob
// package.
// func _CopyMap(val map[string]interface{}) (map[string]interface{}, error) {
// 	var cp bytes.Buffer
// 	enc := gob.NewEncoder(&cp)
// 	if err := enc.Encode(val); err != nil {
// 		return nil, err
// 	}
// 	dec := gob.NewDecoder(&cp)
// 	var result map[string]interface{}
// 	if err := dec.Decode(&result); err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }
