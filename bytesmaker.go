/* This package should be under go/src/bytesmaker to be properly
   imported
*/

package bytesmaker 

import "unsafe"
import "fmt"

/* Converts data of types:
   int, int32, int64, byte and string 
   to a byte array. Least significant byte is at return 
   arrays index 0.

   If multiple arguments, all arguments are converted 
   in order to one combined slice of bytes, i.e:
                                |      int      |string|byte|
   Bytes("ABC", 1, byte(2)) == [1 0 0 0 0 0 0 0 65 66 67 2]

   Examples:
   x := Bytes("Hello world!")
   x := Bytes(2345)
   x := Bytes(int64(1e6), int32(2), "Hello", 4)
*/
func Bytes(data ... interface{}) []byte {

    if len(data) > 1 { // If multiple arguments sent
        res := make([]byte, 0, 8*len(data))
        for _, d := range data {
            res = append(res, Bytes(d) ... )
        }
        return res
    } else {
        data := data[0]
        res := make([]byte, 8)
        switch t := data.(type) {
        case byte:
            x := data.(byte)
            res = []byte{x}
        case int:
            x := data.(int)
            res = make([]byte, unsafe.Sizeof(x))
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
            str := fmt.Sprintf("Bytes(): unexpected type %T\n", t)
            panic(str)
        }
        return res
    }
}