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
#include <fcntl.h>  // fcntl
#include <netinet/in.h>
#include <stdio.h>   // fprintf
#include <stdlib.h>  // atoi, exit
#include <string.h>  // memcpy
#include <sys/epoll.h>  // epoll_create1, epoll_ctl, epoll_wait
#include <sys/errno.h>
#include <sys/socket.h>  // socket
#include <sys/types.h>
#include <unistd.h>  // getopt, sleep

#include <limits>  // numeric_limits
#include <string>

#include "net_util.h"

using std::numeric_limits;
using std::string;

// Number of events to handle in each iteration.
const int kNumEvents = 64;

// Timeout for epoll_wait(2).
const int kEpollWaitTimeoutMs = 200;

// Send back a response after each request.
const char kResponse[] = "I got your note.\n";

// Loop until we get this message.
const char kQuit[] = "quit";

void Usage() {
  fprintf(stderr,
          "Usage: server [-p port] [-s seconds]\n"
          " -p\tPort to listen on.\n"
          " -s\tSleep this many seconds to simulate work.\n"
         );
  exit(EXIT_FAILURE);
}

int main(int argc, char* argv[]) {
  int ch;
  uint16_t opt_port = 8080;
  int opt_sleep = 0;

  while ((ch = getopt(argc, argv, "p:")) != -1) {
    switch (ch) {
    case 'p':
      {
        int p = atoi(optarg);
        if (p < numeric_limits<uint16_t>::min() ||
            p > numeric_limits<uint16_t>::max()) {
          Usage();
        }
        opt_port = static_cast<uint16_t>(p);
      }
      break;
    case 's':
      opt_sleep = atoi(optarg);
      if (opt_sleep < 0 || opt_sleep > 86400) {
        fprintf(stderr, "-s seconds should be between 0 and 86,400 seconds.\n");
        Usage();
      }
      break;
    default:
      Usage();
      break;
    }
  }
  argc -= optind;
  argv += optind;

  int listen_fd = socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP);
  if (listen_fd == -1) {
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
  socket_address.sin6_port = htons(opt_port);

  if (bind(listen_fd,
           reinterpret_cast<struct sockaddr*>(&socket_address),
           sizeof(socket_address)) == -1) {
    perror("bind failed");
    return EXIT_FAILURE;
  }

  if (!SetFdNonBlocking(listen_fd)) {
    return EXIT_FAILURE;
  }

  // TODO: Is backlog ignored under Linux?
  if (listen(listen_fd, 10) == -1) {
    perror("listen failed");
    return EXIT_FAILURE;
  }

  // TODO: Print out /proc/sys/fs/epoll/max_user_watches
  int epoll_fd = epoll_create1(0);
  if (epoll_fd == -1) {
    perror("epoll_create1 failed");
    return EXIT_FAILURE;
  }

  struct epoll_event listen_event;
  listen_event.data.fd = listen_fd;
  listen_event.events = EPOLLIN | EPOLLET;
  if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, listen_fd, &listen_event) == -1) {
    perror("epoll_ctl(listen_fd)");
    return EXIT_FAILURE;
  }

  struct epoll_event events[kNumEvents];

  while (true) {
    int num_events = epoll_wait(epoll_fd, events,
                                kNumEvents, kEpollWaitTimeoutMs);
    if (num_events == -1) {
      if (errno == EINTR) {
        continue;
      } else {
        perror("epoll_wait failed");
        return EXIT_FAILURE;
      }
    }
    for (int i = 0; i < num_events; ++i) {
      if (events[i].events & (EPOLLERR | EPOLLHUP)) {
        // We saw an error or HUP on the fd.
        fprintf(stderr, "fd closed unexpectedly");
        close(events[i].data.fd);
        continue;
      } else if (events[i].data.fd == listen_fd &&
                 events[i].events & EPOLLIN) {
        // Accept from the listening socket.
        struct sockaddr_in6 peer_address;
        socklen_t peer_length;
        int peer_fd = accept(
            listen_fd,
            reinterpret_cast<struct sockaddr*>(&peer_address),
            &peer_length);
        if (peer_fd == -1) {
          perror("accept failed");
          continue;
        }

        char address_string[INET6_ADDRSTRLEN];
        if (inet_ntop(AF_INET6, &peer_address.sin6_addr,
                      address_string, sizeof(address_string)) != NULL) {
          printf("Client address is %s\n", address_string);
          printf("Client port is %d\n", ntohs(peer_address.sin6_port));
        }

        if (!SetFdNonBlocking(peer_fd)) {
          close(peer_fd);
          continue;
        }

        struct epoll_event event;
        event.data.fd = peer_fd;
        event.events = EPOLLIN | EPOLLET;  // TODO: Add EPOLLOUT
        if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, peer_fd, &event) == -1) {
          perror("epoll_ctl(peer_fd)");
          close(peer_fd);
          continue;
        }
      } else {
        // Respond to a client.
        int peer_fd = events[i].data.fd;
        bool done = false;

        while (true) {
          char buf[4096];
          ssize_t n = read(peer_fd, buf, sizeof(buf));
          if (n == -1) {
            // TODO: Handle EINTR properly.
            if (errno != EAGAIN) {
              perror("read failed");
              done = true;
            }
            break;
          } else if (n == 0) {
            done = true;
            break;
          }

          if (opt_sleep > 0) {
            sleep(opt_sleep); // Simulate computation.
          }

          if (write(STDOUT_FILENO, buf, n) == -1) {
            perror("write failed");
            return EXIT_FAILURE;
          }
          printf("\n-- Received %ld bytes from %d\n", n, peer_fd);

          if (strncmp(kQuit, buf, sizeof(kQuit)-1) == 0) {
            printf("Quitting...\n");
            return EXIT_SUCCESS;
          }
        }

        if (done) {
          // TODO: Use EPOLLOUT.
          write(peer_fd, kResponse, sizeof(kResponse));
          close(peer_fd);  // Closing the fd removes it from epoll.
        }
      }
    }
  }

  close(epoll_fd);
  close(listen_fd);

  return EXIT_SUCCESS;
}
