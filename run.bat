@echo off 
setlocal enableextensions 

::cd proxy
::call docker build -t proxy_image .
::cd ..

::cd redis
::call docker build -t redis_image .
::cd ..

::cd helm
::call helm delete redis-sentinel-proxy .
::call helm template .
::call helm install redis-sentinel-proxy .
::cd ..

call go build -o main.exe .
call main.exe

::cd redis
::call go test -v


endlocal 
