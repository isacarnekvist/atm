# Protocol suggestion
## Order of communication
1. Connection started
2. Server states if something should be updated or not,
  this is of course followed by un update if there is one
3. Login request
4. If login request accepted, do regular customer mode 
   communication
5. Logout
6. Start over from step 2

## Package description
Package size is always 10 bytes except for update data
transmissions.
### Customer related
Suggestion for splitting package in to byte parts:

1. Action "op-code" (client) OR Answer accept/decline (server)
2. For future use?
3. Bytes 3-10: 64 bit integer argument i.e. balance, withdrawal amount,
      user card number, password...

### Update related

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
        # -Update statement-   #
        #     yes, update      #->
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


