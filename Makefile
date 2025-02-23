TMP_BUILDDIR := release
ARCH := tar

make-all: 
	mkdir ${TMP_BUILDDIR}
	openssl req -x509 -newkey rsa:4096 -keyout ${TMP_BUILDDIR}/key.pem -out ${TMP_BUILDDIR}/cert.pem -days 365 -nodes -subj "/C=US/ST=California/L=San Francisco/O=MyCompany/OU=IT/CN=example.com"
	$(MAKE) make-server
	$(MAKE) make-client
	

make-server: 
	go build -o ${TMP_BUILDDIR}/ server.go
	GOOS=windows GOARCH=amd64 go build -o ${TMP_BUILDDIR}/server_win.exe server.go
	GOOS=linux GOARCH=amd64 go build -o ${TMP_BUILDDIR}/server_nix server.go
	GOOS=darwin GOARCH=arm64 go build -o ${TMP_BUILDDIR}/server_macos server.go
	chmod +x ${TMP_BUILDDIR}/server*
	
make-client:
	go build -o ${TMP_BUILDDIR}/ client.go
	GOOS=windows GOARCH=amd64 go build -o ${TMP_BUILDDIR}/client_win.exe client.go
	GOOS=linux GOARCH=amd64 go build -o ${TMP_BUILDDIR}/client_nix client.go
	GOOS=darwin GOARCH=arm64 go build -o ${TMP_BUILDDIR}/client_macos client.go
	chmod +x ${TMP_BUILDDIR}/client*

