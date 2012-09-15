// A TCP server for exercising the socket API.
//
// This can be run as:
//
//   server -p 8080
//
// and tested with:
//
//   nc localhost 8080
//   nc ::1 8080
//
// Build with:
//
//   make server

#include <arpa/inet.h>  // htons
#include <netinet/in.h>
#include <stdio.h>   // fprintf
#include <stdlib.h>  // atoi, exit
#include <string.h>  // memcpy
#include <sys/errno.h>
#include <sys/socket.h>  // socket
#include <sys/types.h>
#include <unistd.h>  // getopt, sleep

#include <limits>  // numeric_limits
#include <string>

using std::numeric_limits;
using std::string;

void Usage() {
  fprintf(stderr,
          "Usage: server [-p port]\n"
         );
  exit(EXIT_FAILURE);
}

int main(int argc, char* argv[]) {
  int ch;
  uint16_t port = 8080;

  while ((ch = getopt(argc, argv, "p:")) != -1) {
    switch (ch) {
    case 'p':
      {
        int p = atoi(optarg);
        if (p < numeric_limits<uint16_t>::min() ||
            p > numeric_limits<uint16_t>::max()) {
          Usage();
        }
        port = static_cast<uint16_t>(p);
      }
      break;
    default:
      Usage();
      break;
    }
  }
  argc -= optind;
  argv += optind;

  int s = socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP);
  if (s == -1) {
    perror("socket create failed");
    return EXIT_FAILURE;
  }

  // Setup the socket address and bind.  AF_INET6 and in6addr_any (::1)
  // means this will bind to both IPv4 and IPv6.  Mapped IPv4 addresses
  // and more are described in RFC 2553.
  struct sockaddr_in6 socket_address;
  memset(&socket_address, 0, sizeof(socket_address));
  socket_address.sin6_family = AF_INET6;
  socket_address.sin6_addr = in6addr_any;
  socket_address.sin6_port = htons(port);

  if (bind(s, reinterpret_cast<struct sockaddr*>(&socket_address),
           sizeof(socket_address)) == -1) {
    perror("bind failed");
    return EXIT_FAILURE;
  }

  // TODO: Is backlog ignored under Linux?
  if (listen(s, 10) == -1) {
    perror("listen failed");
    return EXIT_FAILURE;
  }

  while (true) {
    // Loop until we get "quit".
    struct sockaddr_in6 peer_address;
    socklen_t peer_length;
    int peer_sock = accept(
        s, reinterpret_cast<struct sockaddr*>(&peer_address), &peer_length);
    if (peer_sock == -1) {
      perror("accept failed");
      return EXIT_FAILURE;
    }

    char buf[4096];
    memset(buf, 0, sizeof(buf));
    size_t n = read(peer_sock, buf, sizeof(buf));
    if (n == -1) {
      perror("read failed");
      return EXIT_FAILURE;
    }
    printf("Received message:\n[%s]\n", buf);

    sleep(2); // Simulate computation.

    write(peer_sock, "I got your note.\n", 16);
    close(peer_sock);

    if (strncmp("quit", buf, 4) == 0) {
      break;
    }
  }

  close(s);
 
  return EXIT_SUCCESS;
}
