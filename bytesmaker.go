/* This package should be under go/src/bytesmaker to be properly
   imported
*/

package bytesmaker 

/* Converts data of types int, int32, int64 and string 
   to a byte array. Least significant byte is at return 
   arrays index 0. 
   Examples:
   x := data2bytes("Hello world!")
   x := data2bytes(2345) 
*/
func Bytes(data interface{}) []byte {
    res := make([]byte, 8)
    switch data.(type) {
    case int:
        res = make([]byte, 4)
        x := data.(int)
        for i, _ := range res {
            res[i] = byte(x & 0xff)
            x = x >> 8
        }
    case int32:
        res = make([]byte, 4)
        x := data.(int32)
        for i, _ := range res {
            res[i] = byte(x & 0xff)
            x = x >> 8
        }
    case int64:
        res = make([]byte, 8)
        x := data.(int64)
        for i, _ := range res {
            res[i] = byte(x & 0xff)
            x = x >> 8
        }

    case string:
        x := data.(string)
        res = make([]byte, len(x))
        for i, _ := range res {
            res[i] = byte(x[i])
        }

    default:
        panic("unexpected type")
    }

    return res
}