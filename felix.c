#include <stdio.h>

int main() {
    char x[] = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
    
    /* Byte "2 till 5" skrivs ut */
    printf("{");
    for (int i = 1; i < 5; i++) {
        printf("%d,", x[i]);
    }
    printf("}\n");

    long *y = (long *) &x[1];
    *y = (*y >> 32);

    /* Samma sak igen */
    printf("{");
    for (int i = 1; i < 5; i++) {
        printf("%d,", x[i]);
    }
    printf("}\n");

    /*          | innan |
        [9 8 7 6 5 4 3 2]
               >> 32
                | efter |
        [0 0 0 0 9 8 7 6] (5 4 3 2)
     */

    return 0;
}