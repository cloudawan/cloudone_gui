go build
mv kubernetes_management_gui docker/kubernetes_management_gui/
find ! -wholename './docker/*' ! -wholename './docker' ! -wholename '.' -exec rm -rf {} +
mv docker/version version
mv docker/environment environment
