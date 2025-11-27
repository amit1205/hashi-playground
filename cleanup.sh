echo "Cleaning up containers ..."
echo "=========================="

docker ps -a | grep -E "orders-api|reporting-worker" | awk '{print $1}' | xargs docker rm -f

echo "Cleaning up images .."
echo "=========================="
docker rmi orders-api:dev reporting-worker:dev
