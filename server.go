package main

import (
    "bufio"
    "bytesmaker"
    "errors"
    "fmt"
    "net"
    "os"
    "strconv"
    "atm/updater"
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
    server_set_temp_pwd_error   = 0x28 // 40
    server_set_balance          = 0x29 // 41
    server_set_withd_prompt     = 0x2a // 42
    server_set_withd_success    = 0x2b // 43
    server_set_logout           = 0x2c // 44
    server_set_main             = 0x2d // 45
    server_no_updates           = 0x2f // 47
)

/* Data about user stored in this struct */
type User struct {
    id          int
    password    int64 
    balance     int64
    temp_code   []int       /* A list of single use numbers to use when withdrawing */
    temp_index  int         /* The index says what single use code should be used next */
}

/* Info about users is stored in this dictionary with user number
 * as key
 */
var user_db map[int64]User

/*
 * Starts with zero and gets incremented when updates have been added
 * Each go-routine keeps track of their clients current version number,
 * this always starts with zero since clients don't save state when
 * reboot
 */
var latest_client_version int
var update_handler *updater.Updater

/* Listens for incoming connection requests and start
 * one go routine per connection.
 */
func main() {

    /* Initialize */
    print("Starting server\n")
    init_user_db()
    update_handler = updater.NewUpdater()
    latest_client_version = 0

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
            update_handler.Update_menu()
            latest_client_version++
        case 9:
            quit <- 1
            return
        default:
            fmt.Printf("Not a valid choice \n")
        }
    }
}

/* 
 * Listens for connections and starts separate
 * go routines to handle each one.
 */
func start_listening() {

    var port string
    if len(os.Args) == 2 {
        port = os.Args[1]
    } else {
        port = "8080"
    }

    ln, err := net.Listen("tcp", ":" + port)
    if err != nil {
            fmt.Printf("Couldn't start server \n")
            panic(err)
        } else {
            fmt.Printf("Server started successfully \n")
        }

    /* Start connection with all requests */
    for {
        c, _ := ln.Accept()
        go handleConnection(c)
    }
}

/* 
 * This handles that a connection goes between states
 * in a correct manner and returns if connection was
 * lost
 */
func handleConnection(c net.Conn) {
    print("Connection started\n")
    client_version := 0
    for {
        fmt.Printf("Latest version: %d \n", latest_client_version)
        fmt.Printf("Client version: %d \n", client_version)
        err := state_updates(c, &client_version)
        if err != nil { break }

        user, err2 := state_login(c)
        if err2 != nil { break }

        /* Do no go to user state if client requested logout */
        if user.id != 0 {   
            err = state_user(user, c)
            if err != nil { break }
        }
    }
    c.Close()
    print("Connection closed\n")
}


/***************************************/
/*                                     */
/*      State management methods       */
/*                                     */
/***************************************/


/* This handles all communication in the state UPDATES */
func state_updates(c net.Conn, client_version *int) error {
    fmt.Printf("Entered UPDATES state \n")

    if *client_version < latest_client_version {
        update_handler.UpdateClient(c)
        *client_version = latest_client_version
    } else {
        send_ten( server_no_updates, 0, c )
    }

    /* 
       Since not listening to client, we cannot
       discover errors from the client in this
       state and therefore return nil at this
       point
     */
    return nil
}


/*
 * Listens for login requests from client
 * Returns: 
 *  - A confirmed user and nil if successful
 *  - Empty User{} and error if connection lost
 *  - Empty User{} and nil if logout was requested
 */
func state_login(c net.Conn) (User, error) {
    
    fmt.Printf("Entering LOGIN state \n")

    /* Initialize local variables */
    valid_id_sent := false
    user          := User{}

    for {
        op, val, err := read_and_decode(c)
        if err != nil { return User{}, err }                    /* Connection was probably closed */

        switch op {


        /* Client sent id */
        case login_number:

            user_temp, user_exists := user_db[val]
            if user_exists {
                send_accept(c)
                user = user_temp;
                valid_id_sent = true;
            } else {
                send_decline(c)
            }


        /* Client sent password */
        case login_pwd:                                         

            if valid_id_sent { 
                if user.password == val {
                    send_accept(c)
                    return user, nil
                } else {
                    send_decline(c)
                }
            } else {
                return User{}, errors.New("Unexpected password before successful login")    
            }


        /* Client sent logout request */
        case user_logout:

            return User{}, nil

        default:

            send_error(c)
            return User{}, errors.New("Unexpected op code in login")

        }
    }
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
                user_db[int64(user.id)] = user                  /* Changed user-state needs to be written back */
                send_ten( server_accept, 0, c )
            } else {
                send_ten( server_decline, 0, c )
            }
        case user_logout:
            fmt.Printf("User logged out \n")
            return nil
        default:
            send_error(c)                                       /* Respond with error */
            fmt.Printf("Client sent unexpected op code \n")
            return errors.New("Unexpected op code")             /* Close connection */
        }
    }

    return nil
}



/***************************************/
/*                                     */
/*            Initializers             */
/*                                     */
/***************************************/

/*
 * Creates a database with some users
 */
func init_user_db() {
    /* Create database */
    user_db = make(map[int64]User)

    /* Add some users */
    user_db[86] = User {
        id:         86,
        password:   1234,
        balance:    1000,
        temp_code:  odd_ints(),
        temp_index: 0,
    }
    user_db[85] = User {
        id:         85,
        password:   1111,
        balance:    2000,
        temp_code:  odd_ints(),
        temp_index: 0,
    }
}

/*
 * Returns a slice of all odd numbers 1 - 99
 * Used for creating single use-codes when
 * withdrawing
 */
func odd_ints() []int {
    x := make([]int, 50)
    for i := 0; i < 50; i += 1 {
        x[i] = i*2 + 1
    }
    return x
}

/***************************************/
/*                                     */
/* A lot of convenience methods follow */
/*                                     */
/***************************************/

/*
 * Reads and returns op-code, 64-bit value, error
 */
func read_and_decode(c net.Conn) (int, int64, error) {
    data := make([]byte, 10)
    _, err := c.Read(data)
    op := bytesmaker.Int(data[0:1])
    val := bytesmaker.Int(data[1:9])
    return op, int64(val), err
}

/* 
 * Scans an unsigned integer from stdin
 * Conventient for menu
 *
 * Returns -1 if input was not digit
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
 * Constructs ten bytes from opcode and value and sends
 * through c
 */
func send_ten(op int, val int64, c net.Conn) {
    data := bytesmaker.Bytes( byte(op), val, byte(0) )
    c.Write(data)
}

/* 
 * Simply sends a server-decline op-code with rest set to zero 
 */
func send_decline(c net.Conn) {
    send_ten( server_decline, 0, c )
}

/* 
 * Simply sends a server-error op-code with rest set to zero 
 */
func send_error(c net.Conn) {
    send_ten( server_error, 0, c )
}

/* 
 * Simply sends a server-accept op-code with rest set to zero 
 */
func send_accept(c net.Conn) {
    send_ten( server_accept, 0, c )
}
