# FELAC™ ATM
This is the worlds most trusted ATM

# Protocol suggestion
DO NOT change existing op-codes since these are now implemented!

## Order of communication
The communication can exist in three states after the communication
is established.

1. UPDATES
    * Packages from server only
    * Client must not send in this state
    * Server must end this state with a "no (more) updates"
2. LOGIN
    1. Client sends user number
    2. Server declines or accepts, if decline start over from previous step  
    3. Client sends user password
    4. Server declines or accepts, if decline start over from previous step
    5. If server accepts, proceed to next state
3. USER
    1. Client sends one package only
    2. If client sent logout, go to UPDATES, otherwise server replies with 
    one package only. The answer package is accept/decline. For balance,
    the 64-bit field holds the answer.
4. Start over from UPDATES state

## Package description
Package size is always 10 bytes except for update data
transmissions.
### States LOGIN and USER
Below is suggestion for splitting package in to byte parts.
When client sends a request, the answer doesn't need a certain
"op-code" in the reply. An answer with the correct format will
be sufficient.

If unexpected op-code is sent, return error op-code and close connection!

- Byte 1: Action "op-code" (client) OR Answer accept/decline (server)
- Option 1 (everything except withdrawal):
    * Bytes 2-9: 64 bit integer argument for balance, user card number and password
- Option 2 (for withdrawal):
    * Bytes 2-5: Single use passcode
    * Bytes 6-9: Amount to withdraw
- Byte 10: Not used, for future use?


<!-- -->
    From client:
    User number         0x00
    Password            0x01
    Balance             0x02
    Withdrawal          0x03    This package should attach single-use code and amount
                                to withdraw in same package, see Option 2 above
    Logout              0x04
 
 
    From server:
    Accept              0x10    The 64 bit field can supply an answer for a question
    Decline             0x11
    Error               0x12    Server should shut down connection because then the 
                                behaviour when protocol isn't followed can simply stay 
                                undefined

### Update related
Suggestion is that all update packages are 10 bytes except for actual
data transmissions. These are always sent from server side only. In 
the real world we should add confirmation from the client with a
checksum.

- Byte 1: Type of update "op-code" (these should be different from op-codes 
          in customer related packages)
- Byte 1-5: Size of first transmission (32 bit int)
- Byte 6-9: Size of second transmission, if the action requires one (32 bit int)
- Byte 10: Not used

The reason the size of two consecutive data transmissions is needed 
is that some actions need two arguments. For example if we want to 
add or change a banner for a certain language, we first need the name
of the language and then the new string two add. So therefore we need
the size of both the language name packages and the following package
with the new banner string.

The actions that need two transmissins are all updates within an already
named language. Argument 1 is the language name, 2 is the new string.

### Different updates
All newlines, if any, should be included in the strings so that design choices
can be freely made in updates.
#### Banner
Max. 80 characters long string.


<!---->
    Args is the number of data sends required after the send containg the 
    op code.

    Description             Code    Args    Example
    ===========================================
    Add language            0x21    1     
    Add/set banner          0x22    2       "Buy stocks in ... bla bla bla. \n"
    Add/set login text      0x23    2       "Please enter user number: \n"
    Add/set passw text      0x24    2       "Please enter password: \n"
    Add/set wrong login     0x25    2       "No such user \n"
    Add/set wrong passw     0x26    2       "Wrong password \n"
    Add/set list passw      0x27    2       "Please enter next password code from list: \n"
    Add/set wrong list pwd  0x28    2       "Wrong list password \n"      # v These were changed!
    Add/set balance text    0x29    2       "Your balance is: \n"
    Add/set withd. amount   0x2a    2       "Enter amount to withdraw: \n"
    Add/set withdrawal text 0x2b    2       "Withdrawal succesful \n"
    Add/set logout          0x2c    2       "You have been logged out \n"
    Add/set main menu       0x2d    2       "Welcome!!! \n"

    No (more) updates       0x2f    0


## Communication example
    ==========================================
    Server                              Client
            ========================
            #  Connection started  #
            ------------------------

            ========================
            #  -Update statement-  #->
            #    no more update    #
            ------------------------

            
            ========================
          <-#       -Login-        #
            #     card number      #
            ------------------------
            ========================
            #    -Login answer-    #->
            #        accept        #
            ------------------------
            ========================
          <-#     -balance inq-    #
            ------------------------
            ========================
            #   -Balance answer-   #->
            #        amount        #
            ------------------------
            ========================
          <-#   -Withdrawal req-   #
            #    single use code   #
            #        amount        #
            ------------------------
            ========================
            #   -Withdrawal ans-   #->
            #    accept/decline    #
            ------------------------
            ========================
          <-#       -Logout-       #
            ------------------------
            
            
            ========================
            #    -Add language-    #->
            #  language name size  #
            ------------------------
            ========================
            #    language name     #->
            ------------------------
            
            ========================
            #  -Set/change banner- #->
            #  language name size  #
            #      banner size     # 
            ------------------------
            ========================
            #    language name     #->
            ------------------------
            ========================
            #        banner        #->
            ------------------------

                       ...

            ========================
            #  -Update statement-  #->
            #    no more update    #
            ------------------------


            ========================
            #  Wait for login req  #
            ------------------------



