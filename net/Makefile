CC=g++
CFLAGS=-I. -Wall -Werror
DEPS=net_util.h
OBJ=net_util.o client.o server.o

%.o: %.c $(DEPS)
	$(CC) -c -o $@ $< $(CFLAGS)

client: client.o net_util.o
	$(CC) -o $@ $^ $(CFLAGS)

server: server.o net_util.o
	$(CC) -o $@ $^ $(CFLAGS)

all: client server
