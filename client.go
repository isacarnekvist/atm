package main 

import (
    "net"
    "os"
    "atm/bytesmaker"
    "strings"
    "bufio"
    "fmt"
    "strconv"
    //"fmt"
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


    //message indexes

    banner                      = 0x00 // 34
    login_prompt                = 0x01 // 35
    passw_prompt                = 0x02  // 36
    userr                       = 0x03 // 37
    wrong_pwd                   = 0x04 // 38
    temp_pwd_prompt             = 0x05 // 39
    temp_pwd_error              = 0x06 // 40
    balance                     = 0x07 // 41
    withd_prompt                = 0x08 // 42
    withd_success               = 0x09 // 43
    logout                      = 0x0a // 44
    main_prompt                 = 0x0b // 45

    langage_length              = 12
    command_start               = 0x22


)


/* Genom att skicka t ex "中文" som key så får man 
   tillbaka den struct som beskriver programmets olika delar. */
var languages map[string][]string

/* Detta är den struct som kommandon läses från
   just nu. */
var current_language [] string


// Reader to be used for input
var reader *bufio.Reader = bufio.NewReader(os.Stdin)

/* MAIN */
func main() {

    // Se till att språk initialiseras och ett start-språk väljs
    currentState := 0 
    init_lang()
    conn, err := net.Dial("tcp", "127.0.0.1:8080")          // Anslut till ip:port
    if err != nil {
        print("Error connecting")
        return
    }
   
    defer conn.Close()
    for  {
        if currentState == 9 {
            break
        }
        stateHandler(&currentState, conn)
    }
}


func stateHandler(state *int, c net.Conn){
    switch *state {
        case 0:
            update_state(state, c)
        case 1:
            login_state_userID(state, c)
        case 2:
            login_state_password(state, c)
        case 3:
            loggedin_state(state, c)
        default :
            *state=9
            update_state(state, c)
    }
}

func update_state(state *int, c net.Conn) {
    for *state == 0 {

        op, value, err := read_and_decode(c)
        size1:=value & 0xffffffff
        size2:= value>>32

        resp1 := make([]byte, size1)
        resp2 := make([]byte, size2)
        switch {
            case err!=nil:
                println("error reading update data")
                *state=9
                return 
            case op == server_no_updates:

                *state=1
                language_coice(state, c)
                return

            case op == server_set_language:
                c.Read(resp1)
                new_lang := make([]string, langage_length)
                languages[string(resp1)]= new_lang
                return

            case op >= server_set_banner && op <= server_set_main:
                c.Read(resp1)
                c.Read(resp2)
                lang := languages[string(resp1)]
                lang[op-command_start-1]=string(resp2)
                return
            default :
                return
        }
    }
}

func language_coice(state *int, c net.Conn){
    println("Choose language / välj språk")
    lang_list:= make([]string, len(languages))
    count:=0
    for language := range languages {
        lang_list[count] =  language
        fmt.Printf("%d : %s \n", count, language)
        count++
    }

    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    lang, err:= strconv.Atoi(line)
    if err != nil{
        *state= 9
        return
    } else {
       current_language = languages[lang_list[lang]]
    }

}

func login_state_userID(state *int, c net.Conn) {
    print(current_language[login_prompt])
   
    // Code to read from prompt, could be replaced with other io code
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    userID, err := strconv.Atoi(line)

    if line=="c"{
        handle_logout(state, c)
        return

    }

    if err != nil{
        println(current_language[userr])
    } else {
        send_ten(login_number, int64(userID), c )
        op, _, err2 := read_and_decode(c)

        if err2 != nil{
            println("error has occured wile sending userID, program will exit")
            *state=9           
        } else {
            switch op {
                case server_accept:
                    *state = 2
                case server_decline:
                    //println(current_language[userr])
            }
        }
    }
}

func login_state_password(state *int, c net.Conn) {
   print(current_language[passw_prompt])
   
    // Code to read from prompt, could be replaced with other io code
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    password, err:= strconv.Atoi(line)

    if line=="c"{
        handle_logout(state, c)
        return
    }

    if err != nil{
                println(current_language[wrong_pwd])
    } else {
        send_ten(login_pwd, int64(password), c)
        op, _, err2 := read_and_decode(c)

        if err2 != nil{
            println("error has occured wile sending userID, program will exit")
            *state=9           
        } else {
            switch {
                case op==server_accept:
                    *state = 3
                    return
                case op==server_decline:
                    print(current_language[wrong_pwd])
                    return
            }
        }
    }
}



func loggedin_state(state *int, c net.Conn) {
    print(current_language[main_prompt])
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    inp, err:= strconv.Atoi(line)
    if err != nil{
        return
    } else { 
    switch inp {
        case 1:
            handle_balance_request(state, c)
            return
        case 2:
            handle_withdrawal(state, c)
            return
        case 3:
            language_coice(state, c)

        case 4:
            handle_logout(state, c)
            return
        default:
            return
    }
    }
}
func handle_logout( state *int, c net.Conn){
            send_ten(user_logout, int64(0), c)
            *state=0 
            print(current_language[logout])
}


func handle_withdrawal(state *int, c net.Conn){

    print(current_language[temp_pwd_prompt])
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    inp, err := strconv.Atoi(line)
    println(len(line))
    println(line)

    print(current_language[withd_prompt])
    line2, _ := reader.ReadString('\n')
    line2 = strings.TrimSpace(line2)
    w_amount, err2:= strconv.Atoi(line2)

  

    if err != nil || err2!=nil || len(line)>2 {
                print(current_language[temp_pwd_error])
    } else {

        toSend:= int64(w_amount)<<32
        toSend = toSend + int64(inp)
        send_ten(user_withdrawal, toSend, c)
        op, _, err3 := read_and_decode(c)
        
        if err3!=nil { 
            *state = 9 
        } else {
            switch op{
            case server_accept:
                println(current_language[withd_success])
                fmt.Printf("%d\n", w_amount)
            case server_decline:
                print(current_language[temp_pwd_error])
            }
        }
    }
}

func handle_balance_request(state *int, c net.Conn){
    send_ten(user_balance, int64(0), c)
            op, value, err2 := read_and_decode(c)
            if err2!=nil{ 
                *state = 9 
            } else {
                switch op{
                case server_accept:
                    println(current_language[balance])
                    fmt.Printf("%d\n", value)
                default:
                    return
                }
            }

}




 func read_and_decode(c net.Conn) (int, int64, error) {
    data := make([]byte, 10)
    _, err := c.Read(data)
    op := bytesmaker.Int(data[0:1])
    val := bytesmaker.Int(data[1:9])
    return op, int64(val), err
}

func send_ten(op int, val int64, c net.Conn) {
    data := bytesmaker.Bytes( byte(op), val, byte(0) )
    c.Write(data)
}



/* Initiera ett grundtillstånd som gäller innan
   någon uppdatering skett. 
func state_handler( state Integer){

}
*/
func init_lang() {
    languages = make(map[string][]string)
    svenska := make([]string, langage_length)
    current_language = svenska
}