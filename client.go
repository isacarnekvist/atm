package main 

import (
    "net"
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
    init_lang()

    conn, err := net.Dial("tcp", "127.0.0.1:8080")          // Anslut till ip:port
    if err != nil {
        print("Error connecting")
        return
    }

    login_screen(conn)

    /*
    // Skicka data
    data := []byte{0x41, 0x42}
    conn.Write(data)

    // Ta emot!
    resp := make([]byte, 8)
    conn.Read(resp)

    // Skriv ut
    for _, d := range resp {
        fmt.Printf("%x", d)
    }
    print("\n")

    // Stäng
    */
    conn.Close()
}

func login_screen(c net.Conn) {
    print(current_language.login_user)
}

/* Initiera ett grundtillstånd som gäller innan
   någon uppdatering skett. */
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


