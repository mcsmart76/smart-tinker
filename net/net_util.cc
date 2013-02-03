#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>

bool SetFdNonBlocking(int fd) {
  int flags = fcntl(fd, F_GETFL);
  if (flags == -1) {
    perror("fcntl(F_GETFL) failed");
    return false;
  }
  if (fcntl(fd, F_SETFL, flags | O_NONBLOCK) == -1) {
    perror("fcntl(F_SETFL, O_NONBLOCK) failed");
    return false;
  }
  return true;
}
