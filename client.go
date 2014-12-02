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

/* Alla kommandon för ett språk lagras i
 * in denna struct. Denna struktur nås 
 * förslagsvis via en map där key är namnet
 * på språket.
 */
type language_commands struct {
    banner      string      // T ex "Investera i den Grekiska banksektorn, ett säkert val!\n"
    login_user  string      // T ex "Välkommen! \nSkriv ditt kortnummer: \n"
    login_pass  string      // T ex "Skriv in ditt lösenord: \n"
    main_menu   string      // T ex "Huvudmeny, gör ett val \n 1: Saldo \n 2: ..."
    balance     string      // T ex "Ditt saldo är:\n"
}

/* Genom att skicka t ex "中文" som key så får man 
   tillbaka den struct som beskriver programmets olika delar. */
var languages map[string]language_commands

/* Detta är den struct som kommandon läses från
   just nu. */
var current_language language_commands

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
    reader := bufio.NewReader(os.Stdin)
    defer conn.Close()
    for  {
        if currentState == 9 {
            return
        }
        stateHandler(currentState)
    }
}


func stateHandler(state Integer){
    switch state {
        case 0:
            update_state()
        case 1:
            login_state()
        case 2:
            loggedin_state()
        default :
            currentState=0
            update_state()
    }
}

func update_state() {
    while ( currentState == 0){
        resp := make([]byte, 10)
        conn.Read(resp)
        if (resp[0]==0x2f) {
            currentState==1
        } 
        /*
        * Add code to perform update.
        */
}

func login_state(c net.Conn) {
    print(current_language.login_user)
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    userID, err:= strconv.Atoi()
    if err != nil{
        println("Invalid userID")
    } else {
        sendData:=make([]byte, 10)
        bytesmaker.Bytes
        sendHandler()
    }
}


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
        login_pass  : "Skriv in ditt lösenord: \n",
        main_menu   : "Huvudmeny, gör ett val \n 1: Saldo \n 2: ...",
        balance     : "Ditt saldo är:\n",
    }
    languages["svenska"] = svenska

    中文 := language_commands {
        banner      : "你的钱才是我么的！\n",
        login_user  : "欢迎！\n输入你的客户号码: \n",
        login_pass  : "输入你的密码: \n",
        main_menu   : "请做选择 \n 1: 多少钱 \n 2: ...",
        balance     : "你的钱: \n",
    }
    languages["中文"] = 中文

    current_language = svenska
}


