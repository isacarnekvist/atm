package main

import (
    "bytesmaker"
    "fmt"
    "io"
    "net"
    "os/exec"
    "time"
)

const (
    login_number                = 0x00 // 0
    login_pwd                   = 0x01 // 1
    user_balance                = 0x02 // 2
    user_withdrawal             = 0x03 // 3
    user_logout                 = 0x04 // 4
    server_accept               = 0x10 // 16
    server_decline              = 0x11 // 17
    server_error                = 0x12 // 18
    server_set_language         = 0x21 // 33  
    server_set_banner           = 0x22 // 34
    server_set_login_prompt     = 0x23 // 35
    server_set_passw_prompt     = 0x24 // 36
    server_set_userr            = 0x25 // 37
    server_set_wrong_pwd        = 0x26 // 38
    server_set_temp_pwd_prompt  = 0x27 // 39
    server_set_balance          = 0x28 // 40
    server_set_withd_prompt     = 0x29 // 41
    server_set_withd_success    = 0x2a // 42
    server_set_logout           = 0x2b // 43
    server_no_updates           = 0x2f // 47
)

const (
    valid_id    = 86
    valid_pwd   = 1234
    test_bal    = 1000
)

var server_pipe io.WriteCloser
const server_port = "8081"

func main() {
    /* Connect */
    cmd := exec.Command("go", "run", "server.go", server_port)
    var pipe_err error
    server_pipe, pipe_err = cmd.StdinPipe()
    if pipe_err != nil {
        fmt.Printf("Couldn't create pipe \n")
        panic(pipe_err) 
    }
    run_err := cmd.Start()
    if run_err != nil { panic(run_err) }

    fmt.Printf("Please allow server to listen \n")
    time.Sleep(time.Second*2)

    /* Connect */
    c, err := net.Dial("tcp", "127.0.0.1:" + server_port)
    if err != nil { close_and_panic(fmt.Sprintf("Connecting failed \n")) }

    testUpdates(c)
    testLogin(c)

    testUser(c)                                   /* Uses single-use code 1 and 2 */

    /* Test concurrency */
    time.Sleep(time.Millisecond * 500)
    test_vals_with_new_conn(test_bal - 150, 5)    /* 5 is the current single-use code */

    testLogout(c)

    fmt.Printf("OK \nAll tests passed \n")
    fmt.Printf("Quitting server \n")
    server_pipe.Write(bytesmaker.Bytes("9\n"))
}

/* Assume no updates at this point NEED TO CHANGE THIS */
func testUpdates(c net.Conn) {
    read_and_decode(c)
}

func testLogin(c net.Conn) {

    send_ten( login_number, 0xDEADBEEF, c)      /* Wrong id */
    op, _, _ := read_and_decode(c)
    if op != server_decline { close_and_panic(fmt.Sprintf("Did not decline a non-valid id \n")) }

    send_ten( login_number, valid_id, c)        /* Correct id */
    op, _, _ = read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept a valid id \n")) }
    
    send_ten( login_pwd, 6666, c)               /* Wrong password */
    op, _, _ = read_and_decode(c)
    if op != server_decline { close_and_panic(fmt.Sprintf("Did not decline a non-valid password \n")) }
    
    send_ten( login_pwd, valid_pwd, c)          /* Correct password */
    op, _, _ = read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept a valid password \n")) }

}

/* Not implemented yet, copied code! */
func testUser(c net.Conn) {

    send_ten( user_balance, 0, c)               /* Request balance */
    op, val, _ := read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept balance inquiry \n")) }
    if val != 1000 { close_and_panic(fmt.Sprintf("Balance is not the expected one, restart server? \n")) }
    
    send_ten( user_balance, 0, c)               /* Request balance again */
    op, val, _ = read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept second balance inquiry \n")) }
    if val != 1000 { close_and_panic(fmt.Sprintf("Balance is not the same when checked again \n")) }

    send_ten( user_withdrawal, 0xDEADBEEF, c)   /* Request withdrawal with wrong code */
    op, val, _ = read_and_decode(c)
    if op != server_decline { close_and_panic(fmt.Sprintf("Did not decline non-valid single-use code \n")) }

    send_ten( user_balance, 0, c)               /* Request balance again */
    op, val, _ = read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept second balance inquiry \n")) }
    if val != 1000 { close_and_panic(fmt.Sprintf("Balance is not the same when checked again \n")) }

    send_ten_2( user_withdrawal, 1, 100, c)     /* Request withdrawal with correct code */
    op, val, _ = read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not decline non-valid single-use code \n")) }

    send_ten_2( user_withdrawal, 3, 50, c)      /* Request withdrawal with correct code */
    op, val, _ = read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not decline non-valid single-use code \n")) }

    send_ten( user_balance, 0, c)               /* Request balance again */
    op, val, _ = read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept second balance inquiry \n")) }
    if val != 850 { close_and_panic(fmt.Sprintf("Balance is not the expected after withdrawal \n")) }

}

func testLogout(c net.Conn) {

    send_ten( user_logout, 0, c)
    op, _, _ := read_and_decode(c)

    /* Next package from server should be an update statement */
    if op & 0xf0 != 0x20 { close_and_panic(fmt.Sprintf("Server did not send update statement after logout \n")) }

}

func test_vals_with_new_conn(valid_bal int, valid_temp int) {
    /* Connect */
    c, err := net.Dial("tcp", "127.0.0.1:" + server_port)
    if err != nil { close_and_panic(fmt.Sprintf("Connecting failed \n")) }

    read_and_decode(c)

    send_ten( login_number, valid_id, c)        /* Correct id */
    op, _, _ := read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept a valid id \n")) }
    
    send_ten( login_pwd, valid_pwd, c)          /* Correct password */
    op, _, _ = read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept a valid password \n")) }

    send_ten( user_balance, 0, c)               /* Request balance */
    op, val, _ := read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept balance inquiry \n")) }
    if val != int64(valid_bal) { 
        close_and_panic(fmt.Sprintf("Customer have too much money on other ATM (∆ = %d)\n", val - int64(valid_bal))) 
    }

    send_ten_2( user_withdrawal, valid_temp, 100, c)   /* Request withdrawal with valid code */
    op, _, _ = read_and_decode(c)
    if op != server_accept { close_and_panic(fmt.Sprintf("Did not accept valid single-use code \n")) }

    c.Close()

}

/* 
 * Used instead of panic so that started servers gets shut
 * down and stops occupying any port
 */
func close_and_panic(err string) {
    fmt.Printf("Shutting down server \n")
    server_pipe.Write(bytesmaker.Bytes("9\n"))
    panic(err)
}

/*
 * Constructs ten bytes from opcode and value and sends
 * through c
 */
func send_ten(op int, val int64, c net.Conn) {
    data := bytesmaker.Bytes( byte(op), val, byte(0) )
    c.Write(data)
}

/*
 * Constructs ten bytes from opcode and value and sends
 * through c
 */
func send_ten_2(op int, val1 int, val2 int, c net.Conn) {
    val := int64( (val2 << 32) | val1 )
    send_ten(op, val, c)
}

/*
 * Reads and returns op-code, value, error
 */
func read_and_decode(c net.Conn) (int, int64, error) {
    data := make([]byte, 10)
    _, err := c.Read(data)
    op := bytesmaker.Int(data[0:1])
    val := bytesmaker.Int(data[1:9])
    return op, int64(val), err
}