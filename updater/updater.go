/* This package handles client software (languages) and related work
 *  - Keeps a database with the latest client software
 *  - Adds to that database when asked to
 *  - Sends updates to a client when asked to
 */

package updater

import (
    "bufio"
    "bytesmaker"
    "fmt"
    "net"
    "os"
    "strconv"
)

/* Data indices */
const (
    set_language         = 1  
    set_banner           = 2
    set_login_prompt     = 3
    set_passw_prompt     = 4
    set_userr            = 5
    set_wrong_pwd        = 6
    set_temp_pwd_prompt  = 7
    set_temp_pwd_error   = 8
    set_balance          = 9
    set_withd_prompt     = 10
    set_withd_success    = 11
    set_logout           = 12
    set_main             = 13
    save_updates         = 99
)

/* Connection op-codes */
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

type Updater struct {
    /* Maps language names to a string array where the location with
     * index given by the constants above contains that string */
    language_db map[string][]string
}

func NewUpdater() *Updater {
    res := &Updater{}
    res.language_db = make(map[string][]string)
    res.init_language_base()

    return res
}

/*
 * Creates an 'out-of-the-box' language database
 */
func (u *Updater) init_language_base() {
    phrases := make([]string, 14)

    phrases[set_main]           = "Välkommen! Gör ett val: \n" + 
                                  "1) Se saldo \n" + 
                                  "2) Gör uttag \n" +
                                  "3) Byt språk \n" +
                                  "4) Logga ut \n"

    phrases[set_banner]         = "Köp grekiska aktier, för en säker pension! \n"
    phrases[set_login_prompt]   = "Skriv in ditt användar-id för att fortsätta: \n"
    phrases[set_userr]          = "Finns inget sådant id \n"
    phrases[set_passw_prompt]   = "Skriv in ditt lösenord, eller c för att logga ut: \n"
    phrases[set_wrong_pwd]      = "Fel lösenord \n"     /* Efter denna skrivs ovanstående väl igen? */

    phrases[set_withd_prompt]   = "Skriv in belopp att ta ut: \n"
    phrases[set_temp_pwd_prompt]= "Skriv in nästa lösenord från listan: \n"
    phrases[set_temp_pwd_error] = "Fel kod \n"        /* Samma här? */
    phrases[set_withd_success]  = "Uttag lyckades \n"

    phrases[set_balance]        = "Ditt saldo är: \n"

    phrases[set_logout]         = "Du är nu utloggad. \n"

    u.language_db["svenska"] = phrases
    
    phrases2 := make([]string, 14)

    phrases2[set_main]           = "您好! 请做选择: \n" + 
                                  "1) 多少钱 \n" + 
                                  "2) 拿钱 \n" +
                                  "3) 另外语言 \n" +
                                  "4) 结束 \n"

    phrases2[set_banner]         = "你应该在我们的银行保存！ \n"
    phrases2[set_login_prompt]   = "输入客户号码： \n"
    phrases2[set_userr]          = "客户号码不对 \n"
    phrases2[set_passw_prompt]   = "输入你的密码, 还是 ‘c’ \n"
    phrases2[set_wrong_pwd]      = "密码不对 \n"     /* Efter denna skrivs ovanstående väl igen? */

    phrases2[set_withd_prompt]   = "您想拿多少钱？ \n"
    phrases2[set_temp_pwd_prompt]= "输入下一个密码: \n"
    phrases2[set_temp_pwd_error] = "不对 \n"        /* Samma här? */
    phrases2[set_withd_success]  = "通过了 \n"

    phrases2[set_balance]        = "您有：\n"

    phrases2[set_logout]         = "通过了 \n"

    u.language_db["中文"] = phrases2
}

/* All the following defines methods that can be invoked on an Updater struct */

/* Client update menu */
func (u *Updater) Update_menu() {
    for {
        fmt.Printf("Please enter digit of choice from below: \n" + 
                    "1) Add/set language \n" +
                    "2) Add/set banner \n" +
                    "3) Add/set 'enter user number' message \n" +
                    "4) Add/set 'enter user password' message \n" +
                    "5) Add/set 'wrong id' question \n" +
                    "6) Add/set 'wrong password' message \n" +
                    "7) Add/set 'enter next temp code' question \n" +
                    "8) Add/set 'wrong temp code' message \n" +
                    "9) Add/set balance message \n" +
                    "10) Add/set withdrawal amount question \n" +
                    "11) Add/set withdrawal successful message \n" +
                    "12) Add/set logged out message \n" +
                    "13) Add/set main menu message \n" +
                    "99) Save and leave update manager \n")
        choice := scan_uint()
        
        if (choice == save_updates) {
            /* Nothing more to add */
            break
        } else if choice == set_language { 
            u.addLanguage()
        } else if (choice >= set_language && choice <= set_main){
            u.addString(choice)
        } else {
            fmt.Printf("Not a valid choice \n")
        }
    }
}

/* Add a new language with the given name to the database */
func (u *Updater) addLanguage() {

    // Prompt for language name
    fmt.Printf("Enter language name: ")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    name := scanner.Text()

    _, lang_exists := u.language_db[name]
    if lang_exists {
        fmt.Printf("Language already in database, nothing was added \n")
    } else {
        fmt.Printf("Added language: %s \n", name)
        u.language_db[name] = make([]string, 14)
    }
}

/*
 * Arguments:
 * str_type: this is one of the constants defined at the top, it says
 *           which one of the client strings should be added
 */
func (u *Updater) addString(str_type int) {

    // Prompt for language name
    fmt.Printf("Enter language name: ")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    lang := scanner.Text()
    l, lang_exists := u.language_db[lang]

    if !lang_exists {
        fmt.Printf("No language named '%s', nothing was added \n", lang)
    } else {

        // Prompt for the new string to add
        data := ""
        fmt.Printf("Enter new (can be multiline) string, end with EOF (Ctrl-D): ")
        for scanner.Scan() {
            data += scanner.Text() + "\n"
        }

        fmt.Printf("Added string to language: %s \n", lang)
        l[str_type] = data

    }
}

/* 
 * Sends entire database to client through supplied connection */
func (u *Updater) UpdateClient(c net.Conn) {
    for language := range u.language_db {

        /* Add language */
        fmt.Printf("Adding language: %s \n", language)
        lang_str_len := len(language)
        send_ten( server_set_language, int64(lang_str_len) , c)
        c.Write(bytesmaker.Bytes(language))

        /* Send all strings */
        fmt.Printf("Sending language data... \n")
        data_strings := u.language_db[language]
        for i, data := range data_strings {
            op_code := int(0x20) + i

            if len(data) > 0 {
                /* Send lengths of two upcoming sends */
                lengths := int64(  len(data) << 32 | lang_str_len  )
                send_ten(op_code, lengths, c )
                fmt.Printf("Sending: 0x%X lang_len: %d data_len: %d\n", op_code, lang_str_len, len(data))

                /* Send language string and data string */
                c.Write(bytesmaker.Bytes(language))
                c.Write(bytesmaker.Bytes(data))
                fmt.Printf("lang: %s\ndata: %s", language, data)
            }
        }
    }

    /* Notify client that updates are done */
    fmt.Printf("Sending no more updates notification \n")
    send_ten( server_no_updates, 0, c )
}

/* Convenience methods */

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