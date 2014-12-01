# FELAC™ ATM
This is the worlds most trusted ATM

# Protocol suggestion
## Order of communication
1. Connection started
2. Server states if something should be updated or not,
  this is of course followed by an update if there is one
3. Login request
4. If login request accepted, do regular customer mode 
   communication
5. Logout
6. Start over from step 2

## Package description
Package size is always 10 bytes except for update data
transmissions.
### Customer related
Below is suggestion for splitting package in to byte parts.
When client has a request, the answer doesn't need a certain
"op-code" in the reply. An answer with the correct format will
be sufficient.

- Byte 1: Action "op-code" (client) OR Answer accept/decline (server)
- Byte 2-9: 64 bit integer argument i.e. balance, withdrawal amount, 
  user card number, password... The reason for 64 bit is that is that
  I have a lot of CA$H!!!
- Byte 10: Not used, for future use?


<!-- -->
    From client:
    User number         0b0000
    Password            0b0001
    Balance             0b0002
    Withdrawal          0b0003
    Logout              0b0004


    From server:
    Accept              0b1000
    Decline             0b1001

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

### Different updates
All newlines, if any, should be included in the strings so that design choices
can be freely made in updates.
#### Banner
Max. 80 characters long string.


<!---->
    Description             Code        Example
    ===========================================
    Add language            0b1002      
    Add/set banner          0b1003      "Buy stocks in ... bla bla bla. \n"
    Add/set login text      0b1004      "Please enter user number: \n"
    Add/set passw text      0b1005      "Please enter password: \n"
    Add/set wrong login     0b1006      "No such user \n"
    Add/set wrong passw     0b1007      "Wrong password \n"
    Add/set list passw      0b1008      "Please enter next password code from list: \n"
    Add/set balance text    0b1009      "Your balance is: \n"
    Add/set withd. amount   0b100a      "Enter amount to withdraw: \n"
    Add/set withdrawal text 0b100b      "Withdrawal succesful \n"
    Add/set logout          0b100c      "You have been logged out \n"

    No (more) updates       0b1111


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
            #        amount        #
            ------------------------
            ========================
          <-#   -Withdrawal req-   #
            #     one-use-code     #
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



