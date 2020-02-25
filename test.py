#!/usr/bin/env python3
#sudo systemctl start redis-server
#sudo systemctl stop redis-sentinel
#sudo systemctl stop redis-server
#import socket
import time
import redis

r = redis.Redis(host='localhost', port=8080)

while True:
    #a = r.execute_command("get","foo")
    a = r.get('foo')
    print(a)
    time.sleep(1)
r.close()

#s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
#s.connect(('localhost', 8080))
#while True:
#    s.sendall(b'Hello, world')
#    time.sleep(1)
#s.close()