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
    server_set_temp_pwd_error   = 0x2c // 44
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

    testUpdates()

    testLogin()

    testUser()

    fmt.Printf("OK \nAll tests passed \n")
    fmt.Printf("Quitting server \n")
    server_pipe.Write(bytesmaker.Bytes("9\n"))
}

/* Assume no updates at this point NEED TO CHANGE THIS? */
func testUpdates() {
    c := newServer();
    op, _, _ := read_and_decode(c)
    if op != server_no_updates { 
        abort("Server did not send 'no more updates' when just started") 
    }
    closeServer();
}

func testLogin() {
    c := newServer()

    /* Recieve no more updates since just started server */
    read_and_decode(c)

    /* Try non-correct and correct id to login */
    send_ten( login_number, 0xDEADBEEF, c )
    op, _, _ := read_and_decode(c)
    if op != server_decline { 
        abort("Server did not decline wrong id")
    }

    send_ten( login_number, valid_id, c )
    op, _, _ = read_and_decode(c)
    if op != server_accept { 
        abort("Server did not accept valid id") 
    }

    /* Log out after submitting id */
    send_ten( user_logout, 0, c )
    op, _, _ = read_and_decode(c)
    if op != server_no_updates { 
        abort("Server did not send 'no more updates' when user logged out") 
    }

    /* Test different passwords */
    send_ten( login_number, valid_id, c )
    read_and_decode(c)

    send_ten( login_pwd , 0xDEADBEEF, c )
    op, _, _ = read_and_decode(c)
    if op != server_decline { 
        abort("Server did not decline non-valid password") 
    }

    /* Log out after submitting wrong password */
    send_ten( user_logout, 0, c )
    op, _, _ = read_and_decode(c)
    if op != server_no_updates { 
        abort("Server did not send 'no more updates' when user logged out") 
    }

    /* Log in with correct id and password */
    send_ten( login_number, valid_id, c )
    send_ten( login_pwd , valid_pwd, c )
    op, _, _ = read_and_decode(c)
    if op != server_accept { 
        abort("Server did not accept valid password") 
    }

    closeServer()
}

func testUser() {
    fmt.Printf("Testing user state \n")
    c := newServer()
    read_and_decode(c)
    send_ten( login_number, valid_id, c )
    send_ten( login_pwd, valid_pwd, c )

    /* TODO */

    closeServer()
}

func abort(err string) {
    closeServer()
    panic(err)
}

/* 
 * Creates a new server to check agains with
 * reset accounts
 */
func newServer() net.Conn {
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
    if err != nil { 
        fmt.Printf("Connecting failed \n")
        server_pipe.Write(bytesmaker.Bytes("9\n"))
        panic(err)
    }

    return c
}

func closeServer() {
    server_pipe.Write(bytesmaker.Bytes("9\n"))
}

/*
 * Constructs ten bytes from opcode and value and sends
 * through c
 */
func send_ten(op int, val int64, c net.Conn) {
    fmt.Printf("Package sent to server:       op: %.2x val64: %d/0x%X int32(1): %d int32(2): %6.d \n", 
                op, val, val, int32(val & 0xffffff), int32((val >> 32) & 0xffffff))
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
    fmt.Printf("Package recieved from server: op: %.2x val64: %d/0x%X int32(1): %d int32(2): %6.d \n", 
                op, val, val, int32(val & 0xffffff), int32((val >> 32) & 0xffffff))
    return op, int64(val), err
}