#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h> 
#include <netinet/in.h>
#include <arpa/inet.h>

		

main()
{
	int			sockfd, newsockfd ; /* Socket descriptors */
	int			clilen;
	struct sockaddr_in	cli_addr, serv_addr;

	int i;
	char buf[100];		
	if ((sockfd = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
		printf("Cannot create socket\n");
		exit(0);
	}

	
	serv_addr.sin_family		= AF_INET;
	serv_addr.sin_addr.s_addr	= INADDR_ANY;
	serv_addr.sin_port		= htons(6000);

	
	if (bind(sockfd, (struct sockaddr *) &serv_addr,
					sizeof(serv_addr)) < 0) {
		printf("Unable to bind local address\n");
		exit(0);
	}

	listen(sockfd, 5); 
	while (1) {

		clilen = sizeof(cli_addr);
		newsockfd = accept(sockfd, (struct sockaddr *) &cli_addr,
					&clilen) ;

		if (newsockfd < 0) {
			printf("Accept error\n");
			exit(0);
		}

	
		if (fork() == 0) {

			close(sockfd);	
			strcpy(buf,"Authentication Done");
            int a = authentication();
            if (a == 0){
			    send(newsockfd, buf, strlen(buf) + 1, 0);

			    //for(i=0; i < 100; i++) buf[i] = '\0';
			    recv(newsockfd, buf, 100, 0);
			    printf("%s\n", buf);

			    close(newsockfd);
			    exit(0);
            }else{
                 send(newsockfd, buf, strlen(buf) + 1, 0);

			    //for(i=0; i < 100; i++) buf[i] = '\0';
			    recv(newsockfd, buf, 100, 0);
			    printf("%s\n", buf);

			    close(newsockfd);
			    exit(0);
            }
		}

		close(newsockfd);
	}
}

int authentication(){
    int a =0;
    printf("Executing the authenticaion function!!!!!");
	system("/home/ubuntu/scripts/R1.sh");//invoke access request
	system("/home/ubuntu/scripts/qr1.sh");// query permission
    return a;
}
	
