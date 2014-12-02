package main

import (
    "bufio"
    "ATM/Graph.zip/lib"
    "errors"
    "fmt"
    "net"
    "os"
    "strconv"
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

/* Data about user stored in this struct */
type User struct {
    password    int64 
    balance     int64
    temp_code   []int       /* A list of single use numbers to use when withdrawing */
    temp_index  int         /* The index says what single use code should be used next */
}

/* Info about users is stored in this dictionary with user number
 * as key
 */
var user_db map[int64]User

/* Listens for incoming connection requests and start
 * one go routine per connection.
 */
func main() {
    /* Initialize */
    print("Starting server\n")
    init_user_db()

    /* Let clients connect */
    go start_listening()

    /* Handle server maintenance */
    quit := make(chan int)
    go server_prompt(quit)
    select {
        case <- quit:
    }
}

/* Server main menu */
func server_prompt(quit chan int) {
    for {
        fmt.Printf("Please enter digit of choice from below: \n" + 
                   "1) Update clients \n" +
                   "9) Quit server \n")
        choice := scan_uint()
        switch choice {
        case 1:
            print("Not implemented at the moment \n")
        case 9:
            quit <- 1
            return
        default:
            fmt.Printf("Not a valid choice \n")
        }
    }
}

/* Scans an unsigned integer from stdin
 * Returns -1 if error
 */
func scan_uint() int {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    choice, err := strconv.Atoi(scanner.Text())
    if err == nil {
           return choice
    } else {
        return -1
    }
}

/* 
 * Listens for connections and starts separate
 * go routines to handle each one.
 */
func start_listening() {
    ln, _ := net.Listen("tcp", ":8080")
    for {
        c, _ := ln.Accept()
        go handleConnection(c)
    }
}

/* This handles that a connections goes between states
 * in a correct manner and returns if connection was
 * lost
 */
func handleConnection(c net.Conn) {
    print("Connection started\n")
    for {
        err := state_updates(c)
        if err != nil { break }

        user, err2 := state_login(c)
        if err2 != nil { break }

        err = state_user(user, c)
        if err != nil { break }
    }
    c.Close()
    print("Connection closed\n")
}

 /* This handles all communication in the state UPDATES */
func state_updates(c net.Conn) error {
    fmt.Printf("Entered UPDATES state \n")
    d := bytesmaker.Bytes(byte(server_no_updates), byte(0), int64(0))
    fmt.Printf("Sent: %d \n", d)
    _, err := c.Write(d)

    return err
}

/*
 * Listens for login requests from client
 * Returns a valid user or error
 */
func state_login(c net.Conn) (User, error) {
    
    var user User
    fmt.Printf("Entering LOGIN state \n")

    /* Start with user id verification */
    for {
        op, id, err := read_and_decode(c)
        if err != nil { return User{}, err }                    /* Connection was probably closed */
        if op != login_number {                                 /* Client sent unexpected op-code */
            send_ten( server_error, 0, c)                       /* Respond with error */
            fmt.Printf("Client sent unexpected op code \n")
            return User{}, errors.New("Unexpected op code")     /* Close connection */
        }
        usertemp, user_exists := user_db[id]
        user = usertemp

        if user_exists {
            fmt.Printf("Client sent valid id \n")
            send_ten( server_accept, 0, c )
            break                                               /* Go on with passcode check */
        } else {
            fmt.Printf("Client sent non-valid id \n")
            send_ten( server_decline, 0, c )
        }

    }

    /* Password verification */
    for {
        op, pass, err := read_and_decode(c)
        if err != nil { return User{}, err }                    /* Connection was probably closed */
        if op != login_pwd {                                    /* Client sent unexpected op-code */
            send_ten( server_error, 0, c )                      /* Respond with error */
            fmt.Printf("Client sent unexpected op code \n")
            return User{}, errors.New("Unexpected op code")     /* Close connection */
        }

        if user.password == pass {
            fmt.Printf("Client sent valid password \n")
            send_ten( server_accept, 0, c )
            break                                               /* Leave LOGIN state */
        } else {
            fmt.Printf("Client sent non-valid password \n")
            send_ten( server_decline, 0, c )
        }
    }

    fmt.Printf("Leaving LOGIN state \n")

    return user, nil
}

/*
 * Listens for user requests
 * Returns error if connection lost or if unexpected
 * status code was sent
 */
func state_user(user User, c net.Conn) error {
    fmt.Printf("Entering USER state \n")

    for {
        op, val, err := read_and_decode(c)
        if err != nil { return nil }             /* Connection was probably lost */
        switch op {
        case user_balance:
            send_ten(server_accept, user.balance, c)
        case user_withdrawal:
            valid_single_use_code := user.temp_code[user.temp_index]
            code := (val & 0xffffffff)
            amount := (val >> 32) & 0xffffffff
            /* Let's accept negative balance! More income for the bank! */
            if code == int64(valid_single_use_code) {
                user.balance -= amount         
                user.temp_index++   
                send_ten( server_accept, 0, c )
            } else {
                send_ten( server_decline, 0, c )
            }
        case user_logout:
            fmt.Printf("User logged out \n")
            return nil
        default:
            send_ten( server_error, 0, c )                      /* Respond with error */
            fmt.Printf("Client sent unexpected op code \n")
            return errors.New("Unexpected op code")             /* Close connection */
        }
    }

    return nil
}

/*
 * Creates a database with some users
 */
func init_user_db() {
    /* Create database */
    user_db = make(map[int64]User)

    /* Add some users */
    user_db[86] = User {
        password:   1234,
        balance:    1000,
        temp_code:  odd_ints(),
        temp_index: 0,
    }
    user_db[85] = User {
        password:   1111,
        balance:    2000,
        temp_code:  odd_ints(),
        temp_index: 0,
    }
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

/*
 * Constructs ten bytes from opcode and value and sends
 * through c
 */
func send_ten(op int, val int64, c net.Conn) {
    data := bytesmaker.Bytes( byte(op), val, byte(0) )
    c.Write(data)
}


/*
 * Returns a slice of all odd numbers 1 - 99
 */
func odd_ints() []int {
    x := make([]int, 50)
    for i := 0; i < 50; i += 1 {
        x[i] = i*2 + 1
    }
    return x
}
