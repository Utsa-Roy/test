#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<time.h>

int main(){
    clock_t t;
    t = clock();
    system("/home/ubuntu/scripts/R1.sh");
    system("/home/ubuntu/scripts/R2.sh");
    system("/home/ubuntu/scripts/qr1.sh");
    system("/home/ubuntu/scripts/qr2.sh");
    t = clock() - t;
    double time_taken = ((double)t)/CLOCKS_PER_SEC;
    printf("time taken== %f ", time_taken);
    return 0;
}
