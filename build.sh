
# Install HG for one of go package to use
sudo apt-get install -y mercurial

# Install Go 1.5.1
wget https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz

# Go bin path
export PATH=$PATH:/usr/local/go/bin

mkdir -p /tmp/go
# Export Go path
export GOPATH=/tmp/go

# Get Kubernetes management gui
go get github.com/cloudawan/kubernetes_management_gui


go build
mv kubernetes_management_gui docker/kubernetes_management_gui/
find ! -wholename './docker/*' ! -wholename './docker' ! -wholename '.' -exec rm -rf {} +
mv docker/version version
mv docker/environment environment

