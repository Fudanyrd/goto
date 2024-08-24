src = cat.go cp.go hexdump.go nslookup.go parrot.go request.go server.go wget.go
exes = cat.exe cp.exe hexdump.exe nslookup.exe parrot.exe request.exe server.exe wget.exe \
	pwd.exe touch.exe ls.exe tree.exe wc.exe telnet.exe hexless.exe \
	ascii.exe bide.exe

CC = go
CFLAGS = build

all: $(exes)

$(exes): %.exe: %.go
	$(CC) $(CFLAGS) $< 

kjk.txt: 
	echo "锟斤拷千斤拷烫烫烫烫烫" > kjk.txt

urls.txt:
	echo "https://go.dev/dl/go1.23.0.linux-amd64.tar.gz" > urls.txt && \
	echo "https://vcg03.cfp.cn/creative/vcg/800/new/VCG211327790854.jpg" >> urls.txt && \
	echo "https://jyywiki.cn/pages/OS/2022/demos/minimal.S" >> urls.txt && \
	echo "https://go.dev/dl/go1.23.0.linux-amd64.tar.gz" >> urls.txt

.PHONY : clean
clean:
	rm *.exe && go clean -cache
