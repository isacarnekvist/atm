package main 

import (
    "net"
    "os"
    "ATM/Graph.zip/lib"
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
    server_set_balance          = 0x28 // 40
    server_set_withd_prompt     = 0x29 // 41
    server_set_withd_success    = 0x2a // 42
    server_set_logout           = 0x2b // 43
    server_no_updates           = 0x2f // 47
)

/* Alla kommandon för ett språk lagras i
 * in denna struct. Denna struktur nås 
 * förslagsvis via en map där key är namnet
 * på språket.
 */
type language_commands struct {
    banner              string      // T ex "Investera i den Grekiska banksektorn, ett säkert val!\n"
    login_user          string      // T ex "Välkommen! \nSkriv ditt kortnummer: \n"
    login_pass          string      // T ex "Skriv in ditt lösenord: \n"
    login_user_accepted string
    login_pass_accepted string
    login_user_declined string
    login_pass_declined string
    main_menu           string       // T ex "Huvudmeny, gör ett val \n 1: Saldo \n 2: ..."
    single_use_code     string
    withdrawal          string
    balance             string      // T ex "Ditt saldo är:\n"
}

/* Genom att skicka t ex "中文" som key så får man 
   tillbaka den struct som beskriver programmets olika delar. */
var languages map[string]language_commands

/* Detta är den struct som kommandon läses från
   just nu. */
var current_language language_commands


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
        resp := make([]byte, 10)
        c.Read(resp)
        if (resp[0]==0x2f) {
            *state=1
        } 
        /*
        * Add code to perform update.
        */
    }
}



func login_state_userID(state *int, c net.Conn) {
    print(current_language.login_user)
   
    // Code to read from prompt, could be replaced with other io code
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    userID, err:= strconv.Atoi(line)

    if err != nil{
        println(current_language.login_user_declined)
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
                    println(current_language.login_user_accepted)
                case server_decline:
                    println(current_language.login_user_declined)
            }
        }
    }
}

func login_state_password(state *int, c net.Conn) {
   print(current_language.login_pass)
   
    // Code to read from prompt, could be replaced with other io code
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    password, err:= strconv.Atoi(line)

    if err != nil{
        println(current_language.login_pass_declined)
    } else {
        send_ten(login_pwd, int64(password), c)
        op, _, err2 := read_and_decode(c)

        if err2 != nil{
            println("error has occured wile sending userID, program will exit")
            *state=9           
        } else {
            switch op {
                case server_accept:
                    *state = 3
                    println(current_language.login_pass_accepted)
                case server_decline:
                    println(current_language.login_pass_declined)
            }
        }
    }
}



func loggedin_state(state *int, c net.Conn) {
    print(current_language.main_menu)
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    inp, err:= strconv.Atoi(line)
    if err != nil{
        println("Invalid choice")
    } else { 
    switch inp {
        case 1:
            handle_balance_request(state, c)
        case 2:
            handle_withdrawal(state, c)
        case 9:
            *state=9 
        default:
            println("Invalid Choice")
    }
    }
}


func handle_withdrawal(state *int, c net.Conn){

    print(current_language.single_use_code)
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    inp, err := strconv.Atoi(line)

    print(current_language.withdrawal)
    line, _ = reader.ReadString('\n')
    line = strings.TrimSpace(line)
    w_amount, err2:= strconv.Atoi(line)

  

    if err != nil || err2!=nil{
        println("Invalid choice")
    } else {

        toSend:= int64(w_amount)<<32
        toSend = toSend + int64(inp)
        send_ten(user_withdrawal, toSend, c)
        op, value, err3 := read_and_decode(c)
        
        if err3!=nil { 
            *state = 9 
        } else {
            switch op{
            case server_accept:
                println(current_language.balance)
                fmt.Printf("%d\n", value)
            case server_decline:
                println("Withdrawal denied")
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
                    println(current_language.balance)
                    fmt.Printf("%d\n", value)
                case server_decline:
                    println("Balance denied")
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
    languages = make(map[string]language_commands)

    svenska := language_commands {
        banner      : "Investera i den Grekiska banksektorn, ett säkert val!\n",
        login_user  : "Välkommen! \nSkriv ditt kortnummer: \n",
        login_user_accepted : "UserID accepted\n",
        login_user_declined : "Invalid userID\n",
        login_pass  : "Skriv in ditt lösenord: \n",
        login_pass_accepted : "Password accepted\n",
        login_pass_declined : "Invalid password",
        main_menu   : "Huvudmeny, gör ett val \n 1: Saldo \n 2: Uttag, 9: Exit\n",
        single_use_code : "Ange din engångskod: \n",
        withdrawal : "Ange summan du vill ta ut \n",
        balance     : "Ditt saldo är:\n",
    }
    languages["svenska"] = svenska

    中文 := language_commands {
        banner      : "你的钱才是我么的！\n",
        login_user  : "欢迎！\n输入你的客户号码: \n",
        login_user_accepted : "UserID accepted\n",
        login_user_declined : "Invalid userID\n",
        login_pass  : "输入你的密码: \n",
        login_pass_accepted : "Password accepted\n",
        login_pass_declined : "Invalid password \n",
        main_menu   : "请做选择 \n 1: 多少钱 \n 2: ...",
        single_use_code: "Ange din fyrasiffriga engångskod: \n",
        balance     : "你的钱: \n",
    }
    languages["中文"] = 中文

    current_language = svenska
}


