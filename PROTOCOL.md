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

1. Action "op-code" (client) OR Answer accept/decline (server)
2. For future use?
3. Bytes 3-10: 64 bit integer argument i.e. balance, withdrawal amount, 
   user card number, password...


<!-- -->
    From client:
    Login       0x0000
    Balance     0x0001
    Withdrawal  0x0002
    Logout      0x0003


    From server:
    Accept      0x1000
    Decline     0x1001

### Update related
Suggestion is that all update packages are 10 bytes except for actual
data transmissions. These are always sent from server side only. In 
the real world we should add confirmation from the client with a
checksum.

1. Type of update "op-code" (these should be different from op-codes 
   in customer related packages)
2. Size of next transmission (one byte is sufficient since max is 80)
3. Size of second transmisson (needed for some actions)

The reason the size of two consecutive data transmissions is needed 
is that some actions need two arguments. For example if we want to 
add or change a banner for a certain language, we first need the name
of the language and then the new string two add. So therefore we need
the size of both the language name package and the following package
with the new banner string.

<!---->
    Add language        0x1002
    Add/set banner      0x1003
    Add/set login text  0x1004
    Add/set ...         0x....
    No (more) updates   0x1111


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


