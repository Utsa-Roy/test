#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <time.h>

main()
{   
    fork();
	clock_t start_t, end_t, total_t;
	start_t = clock();


	int			sockfd ;
	struct sockaddr_in	serv_addr;

	int i;
	char buf[100];
	if ((sockfd = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
		printf("Unable to create socket\n");
		exit(0);
	}

	serv_addr.sin_family		= AF_INET;
	serv_addr.sin_addr.s_addr	= inet_addr("127.0.0.1");
	serv_addr.sin_port		= htons(6000);

	/* With the information specified in serv_addr, the connect()
	   system call establishes a connection with the server process.
	*/
	if ((connect(sockfd, (struct sockaddr *) &serv_addr,
						sizeof(serv_addr))) < 0) {
		printf("Unable to connect to server\n");
		exit(0);
	}

	
	for(i=0; i < 100; i++) buf[i] = '\0';
	recv(sockfd, buf, 100, 0);
	printf("%s\n", buf);

	for(i=0; i < 100; i++) buf[i] = '\0';
	strcpy(buf,"Message from client");
	send(sockfd, buf, 100, 0);
	end_t = clock();
	total_t = (double)(end_t - start_t) / (CLOCKS_PER_SEC/1000);
   	printf("Total time taken by CPU: %f\n", total_t  );

	close(sockfd);


	
}