// A TCP client for exercising the socket API.
//
// This can be tested against a server such as:
//
//   nc -4 -l -k 8080
//
// and run as:
//
//   client -s localhost -p 8080 -m "hello, world"
//
// Build with:
//
//   make client

// TODO: macros for logging error and perror
// TODO: proper errno handling; helper function using strerror_r
// and convert to a string
// TODO: IPv6 support

#include <arpa/inet.h>  // htons
#include <fcntl.h>  // open
#include <netdb.h>  // gethostbyname
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>  // atoi, exit
#include <string.h>  // memcpy
#include <sys/errno.h>
#include <sys/mman.h>  // mmap, madvise
#ifdef __linux__
#include <sys/sendfile.h>  // sendfile
#endif
#include <sys/socket.h>  // socket, sendfile (OSX)
#include <sys/stat.h>  // fstat
#include <sys/types.h>
#ifdef __APPLE__
#include <sys/uio.h>  // sendfile
#endif
#include <unistd.h>  // getopt

#include <iostream>
#include <limits>  // numeric_limits
#include <string>

using std::numeric_limits;
using std::string;

// Size of chunks to write.
static const size_t kWriteSize = 4096;

// setsockopt SO_LINGER timeout in seconds.
static const int kLingerTimeout = 10;

// Connect via TCP to the endpoint given and return a valid
// socket.  Returns -1 and an error string.
int TcpConnect(const string& server, uint16_t port, bool linger,
               string* error_message) {
  int s = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
  if (s == -1) {
    *error_message = "socket create failed: " + string(strerror(errno));
    return -1;
  }

  if (linger) {
    // Lingering is of marginal use.  Normally an application-level
    // acknowledgement is necessary to know for certain that the server
    // has received and processed a client's data.
    // http://developerweb.net/viewtopic.php?id=2982
    struct linger lingeropt;
    lingeropt.l_onoff = 1;
    lingeropt.l_linger = kLingerTimeout;
    if (setsockopt(s, SOL_SOCKET, SO_LINGER,
                   &lingeropt, sizeof(lingeropt)) == -1) {
      *error_message = "setsockopt(SO_LINGER) failed: " + string(strerror(errno));
      return -1;
    }
  }

  // Resolve the server's name into an IP address.
  // TODO: IPv6
  struct hostent* hostent = gethostbyname(server.c_str());
  if (hostent == NULL) {
    *error_message = "could not resolve server name: " + server + ": " +
        string(hstrerror(h_errno));
    return -1;
  }

  // Setup the socket address and connect to the server.
  struct sockaddr_in address;
  memset(&address, 0, sizeof(address));
  address.sin_port = htons(port);  // TODO: parse
  memcpy(&address.sin_addr.s_addr, hostent->h_addr_list[0], hostent->h_length);
  address.sin_family = hostent->h_addrtype;
  if (connect(s, reinterpret_cast<struct sockaddr*>(&address),
              sizeof(address)) == -1) {
    *error_message = "connect failed: " + string(strerror(errno));
    return -1;
  }

  return s;
}

// Writes message to fd.  Partial writes are handled.
// Returns -1 on error and fills in error_message.
int WriteMessage(int fd, string message, string* error_message) {
  char* buffer = const_cast<char*>(message.data());
  size_t bytes_remaining = message.size();
  size_t write_size = kWriteSize;
  do {
    ssize_t bytes_sent = write(fd, buffer, write_size);
    if (bytes_sent < 0) {
      *error_message = "write failed: " + string(strerror(errno));
      return -1;
    }
    if (bytes_sent > bytes_remaining) {
      *error_message = "write returned more than we sent!";
      return -1;
    }
    buffer += bytes_sent;
    bytes_remaining -= bytes_sent;
  } while (bytes_remaining > 0);

  return message.size();
}

size_t FileSizeFromFd(int fd) {
  struct stat stat;
  if (fstat(fd, &stat) == -1) {
    return 0;
  }
  return stat.st_size;
}

void Usage() {
  fprintf(stderr,
          "Usage: client [-lz] [-s host] [-p port] [-m message]\n"
          " -l\tLinger on socket close.\n"
          " -z\tUse sendfile(2) instead of mmap+write.\n"
         );
  exit(EXIT_FAILURE);
}

int main(int argc, char* argv[]) {
  int ch;
  bool linger = false;
  bool use_sendfile = false;
  string filename;
  string message;
  string server;
  uint16_t port = 8080;

  while ((ch = getopt(argc, argv, "f:lm:p:s:z")) != -1) {
    switch (ch) {
    case 'f':
      filename = string(optarg);
      break;
    case 'l':
      linger = true;
      break;
    case 'm':
      message = string(optarg) + string("\n");
      break;
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
    case 's':
      server = string(optarg);
      break;
    case 'z':
      use_sendfile = true;
      break;
    default:
      Usage();
      break;
    }
  }
  argc -= optind;
  argv += optind;

  int fd = -1;
  void* data = NULL;
  size_t file_size = 0;
  if (!filename.empty()) {
    fd = open(filename.c_str(), O_RDONLY);
    if (fd == -1) {
      fprintf(stderr, "%d: %s\n", __LINE__, strerror(errno));
      return EXIT_FAILURE;
    }

    file_size = FileSizeFromFd(fd);

    if (!use_sendfile) {
      // Read message from mmap'ed file.  This keeps us from copying
      // into userspace, but still means an mmap2() syscall followed
      // by write(2) syscalls.  sendfile(2) can do better.
      //
      // Under Linux we can run:
      //
      //  fcntl(fd, F_SETSIG, RT_SIGNAL_LEASE);
      //  fcntl(fd, F_SETLEASE, F_RDLCK);
      //  mmap(...);
      //  .. use fd ..
      //  fcntl(fd, F_SETLEASE, F_UNLCK);
      //  munmap(...);
      //
      // to be notified if another process attempts to truncate the
      // file while we're sending.
      data = mmap(0, file_size, PROT_READ, MAP_PRIVATE | MAP_NORESERVE, fd, 0);
      if (data == MAP_FAILED) {
        fprintf(stderr, "%d: %s\n", __LINE__, strerror(errno));
        close(fd);
        return EXIT_FAILURE;
      }

      if (madvise(data, file_size, MADV_SEQUENTIAL | MADV_WILLNEED) == -1) {
        fprintf(stderr, "%d: %s\n", __LINE__, strerror(errno));
        close(fd);
        munmap(data, file_size);
        return EXIT_FAILURE;
      }
    }
  }

  if ((message.empty() && filename.empty()) || server.empty() || port == 0) {
    Usage();
  }

  string error_message;
  int peer_socket = TcpConnect(server, port, linger, &error_message);
  if (peer_socket == -1) {
      fprintf(stderr, "%d: %s\n", __LINE__, error_message.c_str());
    return EXIT_FAILURE;
  }

  if (filename.empty()) {
    if (WriteMessage(peer_socket, message, &error_message) == -1) {
      fprintf(stderr, "%d: %s\n", __LINE__, error_message.c_str());
      return EXIT_FAILURE;
    }
  } else if (use_sendfile) {
    // sendfile(2) can do zero-copy with the appropriate interface support,
    // but should be at least as good as the mmap(2) version even without it.
    //
    // OSX's implementation of sendfile(2) allows headers and trailers
    // as iovecs.  Linux, unfortunately, does not.  The standard workaround
    // is to use TCP_CORK.
#ifdef __APPLE__
    off_t len = 0;
    if (sendfile(fd, peer_socket, 0, &len, NULL, 0) == -1) {
#else
#ifdef __linux__
    ssize_t len = sendfile(peer_socket, fd, NULL, 0);
    if (len == -1) {
#endif
#endif
      fprintf(stderr, "%d: %s\n", __LINE__, strerror(errno));
      return EXIT_FAILURE;
    }
    int64_t bytes_sent = static_cast<int64_t>(len);
    if (len != file_size) {
      fprintf(stderr, "%d: did not send complete file: %ld vs. %lu",
              __LINE__, bytes_sent, file_size);
      return EXIT_FAILURE;
    }
  } else {
    // TODO: Merge the -m and -f paths.
    char* current = static_cast<char*>(data);
    int64_t total_bytes_sent = 0;
    while (total_bytes_sent < file_size) {
      ssize_t bytes_sent = send(peer_socket, current, kWriteSize, 0);
      if (bytes_sent == -1) {
        fprintf(stderr, "%d: %s\n", __LINE__, strerror(errno));
        return EXIT_FAILURE;
      }
      if (bytes_sent == 0) {
        fprintf(stderr, "%d: 0-byte send?\n", __LINE__);
        return EXIT_FAILURE;
      }
      total_bytes_sent += bytes_sent;
      current += bytes_sent;
    }
  }

  if (fd != -1) {
    close(fd);
  }
  if (data != NULL && !use_sendfile) {
    munmap(data, file_size);
  }

  // Sends FIN.
  shutdown(peer_socket, SHUT_RDWR);

  // Drain any pending data to read.
  while (true) {
    char buffer[4096];
    ssize_t bytes_read = read(peer_socket, buffer, sizeof(buffer));
    if (bytes_read == -1) {
      fprintf(stderr, "%d: %s\n", __LINE__, strerror(errno));
      return EXIT_FAILURE;
    } else if (bytes_read == 0) {
      break;
    }
  }

  close(peer_socket);

  return EXIT_SUCCESS;
}
